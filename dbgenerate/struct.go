package dbgenerate

import (
	"fmt"
	"github.com/myadamtest/adam/utils"
	"os"
	"strings"
)

type structInfo struct {
	Name        string
	TableName   string
	ProjectName string
	PrivateName string
	FieldInfos  []*structFieldInfo
	PrimaryKey  *structFieldInfo //fixme 暂设假设只是单主键
}

type structFieldInfo struct {
	Name        string
	PrivateName string
	Tp          string
	Tag         string
	Comment     string
}

// sql类型转换golang类型
var typeMap = map[string]string{
	"int":       "int32",
	"decimal":   "float64",
	"varchar":   "string",
	"timestamp": "string",
}

func tableConversion2Struct(info *tableInfo) *structInfo {
	// check

	si := &structInfo{}
	si.Name = camelString(info.Name)
	si.TableName = info.Name
	si.Name = strings.ToUpper(si.Name[0:1]) + si.Name[1:]
	si.PrivateName = strings.ToLower(si.Name[0:1]) + si.Name[1:]

	//fixme 暂设
	//si.ProjectName = "github.com/myadamtest/adam/dbgenerate"
	si.ProjectName, _ = utils.GetProjectName()

	si.FieldInfos = make([]*structFieldInfo, len(info.Fields))
	for i, f := range info.Fields {
		t := "string"
		for key, value := range typeMap {
			if strings.Contains(f.Type, key) {
				t = value
				break
			}
		}

		sfi := &structFieldInfo{}
		si.FieldInfos[i] = sfi
		sfi.Name = camelString(f.Field)
		sfi.Name = strings.ToUpper(sfi.Name[0:1]) + sfi.Name[1:]
		sfi.PrivateName = strings.ToLower(sfi.Name[0:1]) + sfi.Name[1:]
		sfi.Tp = t
		sfi.Comment = fmt.Sprintf("//%s", f.Comment)

		if t == "Time" {
			sfi.Tag = fmt.Sprintf("`json:\"%s\" gorm:\"default:'galeone'\"`", camelString(f.Field))
		} else {
			sfi.Tag = fmt.Sprintf("`json:\"%s\"`", camelString(f.Field))
		}

		if f.Key == "PRI" {
			sfi.Tag = fmt.Sprintf("`json:\"%s\" gorm:\"primary_key\"`", camelString(f.Field))
			si.PrimaryKey = sfi
			continue
		}
	}
	return si
}

func generateStruct(info *structInfo) error {
	err := os.Mkdir("./entity", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	ormFilename := fmt.Sprintf("./entity/orm.go")
	structFilename := fmt.Sprintf("./entity/entity_struct.go")

	_, err = os.Stat(ormFilename)
	if err != nil && os.IsNotExist(err) { //初始化ormPage
		err = executeTemplate(defaultPageTemplate, ormFilename, info, false)
		if err != nil {
			return err
		}
	}

	err = executeTemplate(pageTemplate, ormFilename, info, true)
	if err != nil {
		return err
	}

	_, err = os.Stat(structFilename)
	if err != nil && os.IsNotExist(err) { //初始化ormPage
		structFd, err := os.OpenFile(structFilename, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}

		_, err = structFd.WriteString("package entity \nimport \"github.com/myadamtest/adam/grpcservice/pb/pb\"\n")
		if err != nil {
			structFd.Close()
			return err
		}
		structFd.Close()
	}

	err = executeTemplate(structTemplate, structFilename, info, true)
	if err != nil {
		return err
	}

	return nil
}

func camelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

const structTemplate = `
type {{.Name}} struct {
			pb.{{.Name}}
}
`

const pageTemplate = `
func ({{.Name}}) TableName() string {
	return "{{.TableName}}"
}
`

//// {{.Name}}分页查询条件
//type {{.Name}}Query struct {
//	{{.Name}}
//	Page
//}
//
//// {{.Name}}分页查询结果
//type {{.Name}}Page struct {
//	List []*{{.Name}} ` + "`" + `json:"list"` + "`" + `
//	Page
//}

const defaultPageTemplate = "package entity\n"

//"type Page struct {\n" +
//"\tPageNo   int `json:\"pageNo\" binding:\"required\"`   // 当前页码\n" +
//"\tPageSize int `json:\"pageSize\" binding:\"required\"` // 每页条数\n" +
//"\tTotal    int `json:\"total\"`                       // 总条数\n" +
//"}"
