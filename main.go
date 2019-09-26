package main

import (
	"adam/create"
	"adam/generate"
	"fmt"
	"os"
)

const help = `
adam is create new project and generate code tool
Usage:
		adam <command> [arguments]
The commands are:
		create 			crate a new project.example:adam create projectName
		gen 			gen code, default rpc. support list [grpc].example:adam gen [grpc].
`

const (
	templateUrl  = "https://github.com/myadamtest/adam_template_base.git"
	templateName = "adam_template_base"
)

func main() {
	agrs := os.Args
	if len(agrs) <= 1 {
		fmt.Print(help)
		return
	}

	operation := agrs[1]
	switch operation {
	case "create":
		if len(agrs) < 3 {
			fmt.Print(help)
			fmt.Println("project name can't empty")
			return
		}
		create.CreateProject(agrs[2], templateUrl, templateName)
	case "gen":
		generate.GrpcGenerate()
	default:
		fmt.Print(help)
		return
	}
}
