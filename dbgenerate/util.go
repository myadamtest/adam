package dbgenerate

import (
	"os"
	"text/template"
)

//生成模板命令
func executeTemplate(tpl, fullFilename string, info *structInfo, append bool) error {
	return executeTemplateWithFuncs(tpl, fullFilename, info, append, nil)
}

//生成模板命令
func executeTemplateWithFuncs(tpl, fullFilename string, info *structInfo, append bool, funcs template.FuncMap) error {
	flag := os.O_WRONLY | os.O_CREATE
	if append {
		flag = os.O_WRONLY | os.O_APPEND | os.O_CREATE
	}

	fd, err := os.OpenFile(fullFilename, flag, 0644)
	if err != nil {
		return err
	}

	defer fd.Close()

	templateObject := template.New("")
	if funcs != nil {
		templateObject = templateObject.Funcs(funcs)
	}
	templateObject, err = templateObject.Parse(tpl)
	if err != nil {
		return err
	}

	err = templateObject.Execute(fd, info)
	if err != nil {
		return err
	}
	return nil
}

//判断文件是否存在
func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}

//go类型和proto转换
func type2Proto(tp string) string {
	switch tp {
	case "int", "int8", "int16":
		return "int32"
	case "float32", "float64":
		return "float"
	case "time.Time":
		return "string"
	default:
		return tp
	}
}
