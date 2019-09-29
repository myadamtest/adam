package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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

func GetProjectName() (string, error) {
	b, err := ioutil.ReadFile("./go.mod")
	if err != nil {
		return "", err
	}

	firstLine := strings.Split(string(b), "\n")[0]
	return strings.Trim(strings.Trim(firstLine, " ")[7:], " "), nil
}

func FirstToUpper(str string) string {
	if len(str) == 0 {
		return str
	}

	first := str[0:1]
	return strings.ToUpper(first) + str[1:]
}
