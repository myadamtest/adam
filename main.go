package main

import (
	"adam/create"
)

func main() {
	projectName := "ghyt2"

	create.CreateProject(projectName, "https://github.com/myadamtest/adam_template_base.git", "adam_template_base")
}
