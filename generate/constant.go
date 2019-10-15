package generate

const grpcexeUrl = "https://github.com/myadamtest/grpcexe.git"
const grpcexeFilename = "grpcexe"
const protoBaseDir = "./protofile/"

const (
	commonTemplate = `
package grpcservice

import (
	"{{.ProjectName}}/grpcservice/pb/{{.PackageName}}"
	"context"
)

type {{.Name}}Impl struct {}

{{range $k,$v :=.Interfaces}}
	{{if eq $v.RequestParam.Ty 0}}
		{{if eq $v.ResponseParam.Ty 0}}
func (this *{{$.Name}}Impl) {{$v.Name}}(ctx context.Context,in *{{$.PackageName}}.{{$v.RequestParam.Name}}) (*{{$.PackageName}}.{{$v.ResponseParam.Name}}, error) {
	panic("to implement")
	return nil,nil
}
		{{else if eq $v.ResponseParam.Ty 1}}
func (this *{{$.Name}}Impl) {{$v.Name}}(in *{{$.PackageName}}.{{$v.RequestParam.Name}},st {{$.PackageName}}.{{$.Name}}_{{$v.Name}}Server) (error) {
	panic("to implement")
	return nil
}
		{{end}}
	{{else}}
func (this *{{$.Name}}Impl) {{$v.Name}}(st {{$.PackageName}}.{{$.Name}}_{{$v.Name}}Server) (error) {
	panic("to implement")
	return nil
}
	{{end}}
{{end}}
		`
	methodTemplate = `
{{range $k,$v :=.Interfaces}}
	{{if eq $v.RequestParam.Ty 0}}
		{{if eq $v.ResponseParam.Ty 0}}
func (this *{{$.Name}}Impl) {{$v.Name}}(ctx context.Context,in *{{$.PackageName}}.{{$v.RequestParam.Name}}) (*{{$.PackageName}}.{{$v.ResponseParam.Name}}, error) {
	panic("to implement")
	return nil,nil
}
		{{else if eq $v.ResponseParam.Ty 1}}
func (this *{{$.Name}}Impl) {{$v.Name}}(in *{{$.PackageName}}.{{$v.RequestParam.Name}},st {{$.PackageName}}.{{$.Name}}_{{$v.Name}}Server) (error) {
	panic("to implement")
	return nil
}
		{{end}}
	{{else}}
func (this *{{$.Name}}Impl) {{$v.Name}}(st {{$.PackageName}}.{{$.Name}}_{{$v.Name}}Server) (error) {
	panic("to implement")
	return nil
}
	{{end}}
{{end}}
	`
	//	"{{$.ProjectName}}/grpcservice/pb/{{$.FileName}}"
	//{{$.FileName}}.Register{{$.Name}}Server(s, &{{$.Name}}Impl{})
	grpcStartTemplate = `
package grpcservice

import (
	"{{.ProjectName}}/config"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
	{{range $k,$v :=.PackageList}}{{$v}}
	{{end}}
)

func StartGrpc() error {
	s := grpc.NewServer()
	{{range $k,$v :=.ServiceList}}
	{{$v.PackageName}}.Register{{$v.Name}}Server(s, &{{$v.Name}}Impl{})
	{{end}}
	//...

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GetConfig().RpcPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	reflection.Register(s)

	err = s.Serve(lis)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	return nil
}

func timeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func stringToTime(str string) time.Time {
	t,_ := time.ParseInLocation("2006-01-02 15:04:05",str,time.Local)
	return t
}
`
)
