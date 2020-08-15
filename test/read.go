package main

import (
	"fmt"
	"github.com/gogo/protobuf/plugin/testgen"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	plugin "github.com/gogo/protobuf/protoc-gen-gogo/plugin"
	"go/format"
	"io/ioutil"
	"os"
	"strings"
)

func Read() *plugin.CodeGeneratorRequest {
	g := generator.New()
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		g.Error(err, "reading input")
	}

	if err := proto.Unmarshal(data, g.Request); err != nil {
		g.Error(err, "parsing input proto")
	}

	if len(g.Request.FileToGenerate) == 0 {
		g.Fail("no files to generate")
	}
	return g.Request
}

// filenameSuffix replaces the .pb.go at the end of each filename.
func GeneratePlugin(req *plugin.CodeGeneratorRequest, p generator.Plugin, filenameSuffix string) *plugin.CodeGeneratorResponse {
	g := generator.New()
	g.Request = req
	if len(g.Request.FileToGenerate) == 0 {
		g.Fail("no files to generate")
	}

	g.CommandLineParameters(g.Request.GetParameter())

	g.WrapTypes()
	g.SetPackageNames()
	g.BuildTypeNameMap()
	g.GeneratePlugin(p)

	for i := 0; i < len(g.Response.File); i++ {
		g.Response.File[i].Name = proto.String(
			strings.Replace(*g.Response.File[i].Name, ".pb.go", filenameSuffix, -1),
		)
	}
	if err := goformat(g.Response); err != nil {
		g.Error(err)
	}
	return g.Response
}

func goformat(resp *plugin.CodeGeneratorResponse) error {
	for i := 0; i < len(resp.File); i++ {
		formatted, err := format.Source([]byte(resp.File[i].GetContent()))
		if err != nil {
			return fmt.Errorf("go format error: %v", err)
		}
		fmts := string(formatted)
		resp.File[i].Content = &fmts
	}
	return nil
}

func Generate(req *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
	// Begin by allocating a generator. The request and response structures are stored there
	// so we can do error handling easily - the response structure contains the field to
	// report failure.

	for _, f := range req.ProtoFile {
		if *f.Package == "main" {
			fmt.Println("version ", 3, " >>>>>>>>>>>>>>>>>>>>>>>>>")
			fmt.Println(fmt.Sprintf("%v", f.MessageType[0].Field[1].Name))
			fmt.Println("<<<<<<<<<<<<<<<<<")
		}
	}

	g := generator.New()
	g.Request = req

	g.CommandLineParameters(g.Request.GetParameter())

	// Create a wrapped version of the Descriptors and EnumDescriptors that
	// point to the file that defines them.
	g.WrapTypes()

	g.SetPackageNames()
	g.BuildTypeNameMap()

	g.GenerateAllFiles()

	if err := goformat(g.Response); err != nil {
		g.Error(err)
	}

	testReq := proto.Clone(req).(*plugin.CodeGeneratorRequest)

	testResp := GeneratePlugin(testReq, testgen.NewPlugin(), "pb_test.go")

	for i := 0; i < len(testResp.File); i++ {
		if strings.Contains(*testResp.File[i].Content, `//These tests are generated by github.com/gogo/protobuf/plugin/testgen`) {
			g.Response.File = append(g.Response.File, testResp.File[i])
		}
	}

	return g.Response
}

func Write(resp *plugin.CodeGeneratorResponse) {
	g := generator.New()
	// Send back the results.
	data, err := proto.Marshal(resp)
	if err != nil {
		g.Error(err, "failed to marshal output proto")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		g.Error(err, "failed to write output proto")
	}
}
