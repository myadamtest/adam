package create

import (
	"adam/utils"
	"bytes"
	"fmt"
	"github.com/myadamtest/logkit"
	"io/ioutil"
	"os"
	"os/exec"
)

/**
新建项目
判断项目是否存在
clone项目模板
删除模板git信息
重命名模板
替换相关的项目名信息
执行命令 go mod download
*/
func CreateProject(projectName, templateUrl, templateName string) {
	if projectName == "" {
		logkit.Errorf("project name can't empty")
		return
	}

	_, err := os.Stat(fmt.Sprintf("./%s", projectName))
	if !os.IsNotExist(err) {
		logkit.Errorf("project %s already", projectName)
		return
	}

	cmd := exec.Command("git", "clone", templateUrl)
	cmd.Stderr = os.Stdout
	err = cmd.Run()
	if err != nil {
		logkit.Infof("%s", err)
		return
	}

	err = os.RemoveAll(fmt.Sprintf("./%s/.git", templateName))
	if err != nil {
		logkit.Infof("%s", err)
		return
	}

	err = os.Remove(fmt.Sprintf("./%s/go.sum", templateName))
	if err != nil {
		logkit.Infof("%s", err)
		return
	}

	err = os.Rename(fmt.Sprintf("./%s", templateName), fmt.Sprintf("%s", projectName))
	if err != nil {
		logkit.Infof("%s", err)
		return
	}
	//cmd := exec.Command("")

	err = rangeDirReplaceFile(fmt.Sprintf("./%s", projectName), templateName, projectName)
	if err != nil {
		logkit.Infof("%s", err)
		return
	}

	err = utils.GoModDownload(fmt.Sprintf("./%s", projectName))
	if err != nil {
		fmt.Println("create fail!", err)
		return
	}
	fmt.Println(fmt.Sprintf("project [%s] create success", projectName))
}

func rangeDirReplaceFile(dir, old, new string) error {
	rd, err := ioutil.ReadDir(dir)
	if err != nil {
		logkit.Infof("%s", err)
		return err
	}

	for _, fi := range rd {
		if fi.IsDir() {
			err = rangeDirReplaceFile(fmt.Sprintf("%s\\%s", dir, fi.Name()), old, new)
			if err != nil {
				logkit.Infof("%s", err)
			}
		} else {
			content, _ := ioutil.ReadFile(fmt.Sprintf("%s\\%s", dir, fi.Name()))
			content = bytes.Replace(content, []byte(old), []byte(new), -1)
			err = ioutil.WriteFile(fmt.Sprintf("%s\\%s", dir, fi.Name()), content, 0)
			if err != nil {
				logkit.Infof("%s", err)
			}
		}
	}

	return nil
}
