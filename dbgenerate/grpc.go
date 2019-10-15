package dbgenerate

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

func generateGrpcWithImpl(info *structInfo) error {
	err := generatePb(info.TableName)
	if err != nil {
		return err
	}

	if !fileExist("./grpcservice/pb/pb/common.proto") {
		_ = generatePb("common")
	}

	return generateGrpcImpl(info)
}

func generatePb(tableName string) error {
	_ = os.MkdirAll("./grpcservice/pb/pb", os.ModePerm)

	cmd := exec.Command("protoc", fmt.Sprintf("--go_out=plugins=grpc:./grpcservice/pb/pb"), "--proto_path=./protofile/", fmt.Sprintf("%s.proto", tableName))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func generateGrpcImpl(info *structInfo) error {
	filename := fmt.Sprintf("./grpcservice/%s.go", info.TableName)

	_, err := os.Stat(filename)
	if err == nil {
		// 文件已存在
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}

	tplFuns := template.FuncMap{
		"conversionTpToStruct": conversionTpToStructInTemplate,
		"conversionTypeToPb":   conversionTpToPbInTemplate,
	}

	return executeTemplateWithFuncs(grpcImplTemplate, filename, info, false, tplFuns)
}

func conversionTpToPbInTemplate(tp, paramName, fileName string) string {
	switch tp {
	case "int", "int8", "int16":
		return fmt.Sprintf("int32(%s.%s)", paramName, fileName)
	case "float32", "float64":
		return fmt.Sprintf("float(%s.%s)", paramName, fileName)
	case "time.Time":
		return fmt.Sprintf("%s.%s.String()", paramName, fileName)
	default:
		return fmt.Sprintf("%s.%s", paramName, fileName)
	}
}

func conversionTpToStructInTemplate(tp, paramName, fileName string) string {
	switch tp {
	case "int", "int8", "int16":
		return fmt.Sprintf("%s(%s.%s)", tp, paramName, fileName)
	case "float32", "float64":
		return fmt.Sprintf("%s(%s.%s)", tp, paramName, fileName)
	case "time.Time":
		return "time.Now()"
	default:
		return fmt.Sprintf("%s.%s", paramName, fileName)
	}
}

const grpcImplTemplate = `
package grpcservice

import (
	"fmt"
	"{{.ProjectName}}/entity"
	"{{.ProjectName}}/grpcservice/pb/pb"
	"{{.ProjectName}}/service"
	"context"
)

type {{.Name}}ServiceImpl struct {}

func (this *{{.Name}}ServiceImpl) toStruct(dest *pb.{{.Name}}) *entity.{{.Name}}  {
	{{.PrivateName}} := &entity.{{.Name}}{}

	{{range $k,$v :=.FieldInfos}}{{$.PrivateName}}.{{$v.Name}} = {{conversionTpToStruct $v.Tp "dest" $v.Name}}
	{{end}}
	return {{.PrivateName}}
}

func (this *{{.Name}}ServiceImpl) toPb(dest *entity.{{.Name}}) *pb.{{.Name}}  {
	{{.PrivateName}} := &pb.{{.Name}}{}

	{{range $k,$v :=.FieldInfos}}{{$.PrivateName}}.{{$v.Name}} = {{conversionTypeToPb $v.Tp "dest" $v.Name}}
	{{end}}
	return {{.PrivateName}}
}

func (this *{{.Name}}ServiceImpl) arrToStruct(dest []*pb.{{.Name}}) []*entity.{{.Name}} {
	if len(dest) == 0 {
		return make([]*entity.{{.Name}},0)
	}

	{{.PrivateName}}s := make([]*entity.{{.Name}},len(dest))
	for i,v := range dest {
		{{.PrivateName}}s[i] = this.toStruct(v)
	}
	return {{.PrivateName}}s
}

func (this *{{.Name}}ServiceImpl) arrToPb(dest []*entity.{{.Name}}) []*pb.{{.Name}} {
	if len(dest) == 0 {
		return make([]*pb.{{.Name}},0)
	}

	{{.PrivateName}}s := make([]*pb.{{.Name}},len(dest))
	for i,v := range dest {
		{{.PrivateName}}s[i] = this.toPb(v)
	}
	return {{.PrivateName}}s
}

func (this *{{.Name}}ServiceImpl) pageToStruct(dest *pb.{{.Name}}PageRequest) entity.{{.Name}}Query  {
	{{.PrivateName}}Query := entity.{{.Name}}Query{}
	{{.PrivateName}}Query.Page.PageNo = int(dest.Page.PageNo)
	{{.PrivateName}}Query.Page.PageSize = int(dest.Page.PageSize)
	{{.PrivateName}}Query.Page.Total = int(dest.Page.Total)

	{{.PrivateName}}Query.{{.Name}} = *this.toStruct(dest.{{.Name}})

	return {{.PrivateName}}Query
}

func (this *{{.Name}}ServiceImpl) pageToPb(dest *entity.{{.Name}}Page) *pb.{{.Name}}PageResponse  {
	{{.PrivateName}}Page := &pb.{{.Name}}PageResponse{}
	{{.PrivateName}}Page.Page.PageNo = int32(dest.Page.PageNo)
	{{.PrivateName}}Page.Page.PageSize = int32(dest.Page.PageSize)
	{{.PrivateName}}Page.Page.Total = int64(dest.Page.Total)

	{{.PrivateName}}Page.{{.Name}}List = this.arrToPb(dest.List)
	return {{.PrivateName}}Page
}

func (this *{{.Name}}ServiceImpl) Insert(ctx context.Context,in *pb.{{.Name}}) (*pb.{{.Name}}, error) {
	{{.PrivateName}} := this.toStruct(in)
	err := service.{{.Name}}Service.Insert({{.PrivateName}})
	if err!= nil {
		fmt.Println(err)
		return nil,err
	}

	return this.toPb({{.PrivateName}}),nil
}
		
func (this *{{.Name}}ServiceImpl) Update(ctx context.Context,in *pb.{{.Name}}) (*pb.{{.Name}}, error) {
	{{.PrivateName}} := this.toStruct(in)
	err := service.{{.Name}}Service.Update({{.PrivateName}})
	if err!= nil {
		fmt.Println(err)
		return nil,err
	}

	return this.toPb({{.PrivateName}}),nil
}

func (this *{{.Name}}ServiceImpl) Query(ctx context.Context,in *pb.{{.Name}}PkParamRequest) (*pb.{{.Name}}, error) {
	{{.PrivateName}}, err := service.{{.Name}}Service.Query({{conversionTpToStruct .PrimaryKey.Tp "in" .PrimaryKey.Name}})
	if err!=nil {
		return nil,err
	}

	return this.toPb({{.PrivateName}}),nil
}

		
func (this *{{.Name}}ServiceImpl) Delete(ctx context.Context,in *pb.{{.Name}}PkParamRequest) (*pb.{{.Name}}, error) {
	err := service.{{.Name}}Service.Delete({{conversionTpToStruct .PrimaryKey.Tp "in" .PrimaryKey.Name}})
	if err!=nil {
		return nil,err
	}

	return &pb.{{.Name}}{},nil
}

func (this *{{.Name}}ServiceImpl) QueryList(ctx context.Context,in *pb.{{.Name}}) (*pb.{{.Name}}Array, error) {
	list,err := service.{{.Name}}Service.QueryList(*this.toStruct(in))
	if err!=nil {
		return nil,err
	}

	return &pb.{{.Name}}Array{
			{{.Name}}List:this.arrToPb(list),
			},nil
}

func (this *{{.Name}}ServiceImpl) QueryPage(ctx context.Context,in *pb.{{.Name}}PageRequest) (*pb.{{.Name}}PageResponse, error) {
	result,err := service.{{.Name}}Service.QueryPage(this.pageToStruct(in))
	if err!= nil {
		return nil,err
	}
	return this.pageToPb(result),nil
}
`
