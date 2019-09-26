package utils

import (
	"fmt"
	"os"
	"os/exec"
)

/**
dir 执行目录
*/
func GoModDownload(dir string) error {
	cmd := exec.Command("go", "mod", "download")
	cmd.Stderr = os.Stdout
	if dir != "" {
		cmd.Dir = fmt.Sprintf("./%s", dir)
	}
	return cmd.Run()
}
