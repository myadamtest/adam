package main

import (
	"fmt"
	"github.com/myadamtest/adam/create"
	"github.com/myadamtest/adam/dbgenerate"
	"github.com/myadamtest/adam/generate"
	"os"
)

const help = `
adam is create new project and generate code tool
Usage:
		adam <command> [arguments]
The commands are:
		create 			crate a new project.example:adam create projectName
		gen 			gen code, default rpc. support list [grpc].example:adam gen [grpc].
The command gen:
		grpc			generate code by proto file.
		db 				generate code by db. example: adam gen db "dbuser:password@tcp(ip:port)/dbname"
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
		if len(agrs) >= 3 && agrs[2] == "db" { // 生成数据库
			if len(agrs) < 4 {
				fmt.Print(help)
				fmt.Println("need input db addr.")
				return
			}
			_ = dbgenerate.GenCode(agrs[3])
		} else {
			generate.GrpcGenerate()
		}
	default:
		fmt.Print(help)
		return
	}
}
