package generate

const grpcexeUrl = "https://github.com/myadamtest/grpcexe.git"
const grpcexeFilename = "grpcexe"
const protoBaseDir = "./protofile/"

const (
	commonTemplate = `
package grpcservice

import (
	"{{.ProjectName}}/grpcservice/pb/{{.FileName}}"
	"context"
)

type {{.Name}}Impl struct {}

{{range $k,$v :=.Interfaces}}
	{{if eq $v.RequestParam.Ty 0}}
		{{if eq $v.ResponseParam.Ty 0}}
func (this *{{$.Name}}Impl) {{$v.Name}}(ctx context.Context,in *{{$.FileName}}.{{$v.RequestParam.Name}}) (*{{$.FileName}}.{{$v.ResponseParam.Name}}, error) {
	panic("to implement")
	return nil,nil
}
		{{else if eq $v.ResponseParam.Ty 1}}
func (this *{{$.Name}}Impl) {{$v.Name}}(in *{{$.FileName}}.{{$v.RequestParam.Name}},st {{$.FileName}}.{{$.Name}}_{{$v.Name}}Server) (error) {
	panic("to implement")
	return nil
}
		{{end}}
	{{else}}
func (this *{{$.Name}}Impl) {{$v.Name}}(st {{$.FileName}}.{{$.Name}}_{{$v.Name}}Server) (error) {
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
func (this *{{$.Name}}Impl) {{$v.Name}}(ctx context.Context,in *{{$.FileName}}.{{$v.RequestParam.Name}}) (*{{$.FileName}}.{{$v.ResponseParam.Name}}, error) {
	panic("to implement")
	return nil,nil
}
		{{else if eq $v.ResponseParam.Ty 1}}
func (this *{{$.Name}}Impl) {{$v.Name}}(in *{{$.FileName}}.{{$v.RequestParam.Name}},st {{$.FileName}}.{{$.Name}}_{{$v.Name}}Server) (error) {
	panic("to implement")
	return nil
}
		{{end}}
	{{else}}
func (this *{{$.Name}}Impl) {{$v.Name}}(st {{$.FileName}}.{{$.Name}}_{{$v.Name}}Server) (error) {
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
	{{range $k,$v :=.ServiceList}}
	"{{$v.ProjectName}}/grpcservice/pb/{{$v.FileName}}"
	{{end}}
)

func StartGrpc() error {
	s := grpc.NewServer()
	{{range $k,$v :=.ServiceList}}
	{{$v.FileName}}.Register{{$v.Name}}Server(s, &{{$v.Name}}Impl{})
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
`
)
