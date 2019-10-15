package dbgenerate

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

const (
	daoInitSplit  = "InitDbOrm(addr)"
	daoContentFmt = "%sInitDbOrm(addr)\n\t%sDao = new%sDao()%s"
)

func generateDao(si *structInfo) error {
	err := os.Mkdir("./dao", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	daoFilename := fmt.Sprintf("./dao/%s.go", si.TableName)

	err = executeTemplateWithFuncs(daoTemplate, daoFilename, si, true, template.FuncMap{"setNil": setNil})
	if err != nil {
		return err
	}

	return nil
}

func generateCommonOrm() error {
	err := os.Mkdir("./dao", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	_, err = os.Stat("./dao/orm.go")
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return nil
	}

	return ioutil.WriteFile("./dao/orm.go", []byte(ormDaoTemplate), 0644)
}

func generateIDao(info *structInfo) error {
	filename := "./dao/IDao.go"

	exist, err := daoIsExist(filename)
	if err != nil {
		return err
	}

	if !exist {
		err = executeTemplate(daoInterfaceBase, filename, info, false)
		if err != nil {
			return err
		}
	}

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	content := string(b)
	contentArr := strings.Split(content, daoInitSplit)

	if len(contentArr) != 2 {
		return errors.New("unknown err")
	}
	content = fmt.Sprintf(daoContentFmt, contentArr[0], info.Name, info.Name, contentArr[1])

	err = ioutil.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return err
	}

	err = executeTemplate(daoInterface, filename, info, true)
	if err != nil {
		return err
	}

	return nil
}

func daoIsExist(filename string) (bool, error) {
	_, err := os.Stat(filename)

	if err != nil && !os.IsNotExist(err) {
		return false, err
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return false, err
	}

	str := strings.TrimSpace(string(b))
	if str == "package dao" {
		return false, nil
	}
	return true, nil
}

func setNil(p string) string {
	switch p {
	case "string":
		return ""
	case "int", "int8", "int16", "int32", "int64", "float32", "float64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "0"
	default:
		return "nil"
	}
}

const daoTemplate = `
package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/myadamtest/gobase/logkit"
	"{{.ProjectName}}/entity"
)

type {{.PrivateName}}Dao struct {

}

func new{{.Name}}Dao() I{{.Name}}Dao {
	return &{{.PrivateName}}Dao{}
}

func (dao *{{.PrivateName}}Dao) get{{.Name}}Db() *gorm.DB {
	return baseDb.Model(entity.{{.Name}}{})
}

func (dao *{{.PrivateName}}Dao) Insert({{.PrivateName}} *entity.{{.Name}}) error  {
	err := dao.get{{.Name}}Db().Create({{.PrivateName}}).Error
	if err!=nil {
		logkit.Errorf("insert {{.PrivateName}} err:%s",err)
		return err
	}
	return nil
}

func (dao *{{.PrivateName}}Dao) Update({{.PrivateName}} *entity.{{.Name}}) error  {
	{{if .PrimaryKey}}if isNil({{.PrivateName}}.{{.PrimaryKey.Name}}) {
		return errors.New("delete param can't nil")
	}
	{{end}}
	np := {{.PrivateName}}.{{.PrimaryKey.Name}}
	{{.PrivateName}}.{{.PrimaryKey.Name}} = {{setNil .PrimaryKey.Tp}}
	err := dao.get{{.Name}}Db().Update({{.PrivateName}}).Where("{{.PrimaryKey.ColumnName}}=?",np).Error
	{{.PrivateName}}.{{.PrimaryKey.Name}} = np
	if err!=nil {
		logkit.Errorf("update {{.PrivateName}} err:%s",err)
		return err
	}
	return nil
}

func (dao *{{.PrivateName}}Dao) Query({{.PrimaryKey.PrivateName}} {{.PrimaryKey.Tp}}) (*entity.{{.Name}},error)  {
	result := &entity.{{.Name}}{}
	err := dao.get{{.Name}}Db().Take(result,{{.PrimaryKey.PrivateName}}).Error
	if err!=nil {
		logkit.Errorf("get {{.PrivateName}} err:%s",err)
		return nil,err
	}
	return result,nil
}

func (dao *{{.PrivateName}}Dao) Delete({{.PrimaryKey.PrivateName}} {{.PrimaryKey.Tp}}) error {
	if isNil({{.PrimaryKey.PrivateName}}) {
		return errors.New("delete param can't nil")
	}
	condition := &entity.{{.Name}}{
		{{.PrimaryKey.Name}}:{{.PrimaryKey.PrivateName}},
	}

	err := dao.get{{.Name}}Db().Delete(condition).Error
	if err!=nil {
		logkit.Errorf("del {{.PrivateName}} err:%s",err)
		return err
	}
	return nil
}

func (dao *{{.PrivateName}}Dao) QueryList(filter entity.{{.Name}}) ([]*entity.{{.Name}},error) {
	var result []*entity.{{.Name}}
	err := dao.get{{.Name}}Db().Where(filter).Find(&result).Error
	if err!=nil {
		logkit.Errorf("query list {{.PrivateName}} err:%s",err)
		return nil,err
	}
	return result, nil
}

func (dao *{{.PrivateName}}Dao) QueryPage(q entity.{{.Name}}Query) (*entity.{{.Name}}Page,error)  {
	page := &entity.{{.Name}}Page{}
	page.Page = q.Page

	db := dao.get{{.Name}}Db().Where(q.{{.Name}})
	err := db.Count(&page.Total).Error
	if err!= nil {
		logkit.Errorf("get {{.PrivateName}} total err:%s",err)
		return nil,err
	}

	err = db.Offset((page.PageNo -1) * page.PageSize).Limit(page.PageSize).Find(&page.List).Error
	if err!= nil {
		logkit.Errorf("page query {{.PrivateName}} err:%s",err)
		return nil,err
	}
	return page,nil
}
`

const ormDaoTemplate = `
package dao

import "github.com/jinzhu/gorm"

var baseDb *gorm.DB

func InitDbOrm(addr string)  {
	db,err := gorm.Open("mysql",addr)
	if err!=nil {
		panic(err)
	}

	baseDb = db
}

func isNil(p interface{}) bool {
	if p == nil {
		return true
	}

	switch p.(type) {
	case string:
		return p == ""
	case int,int8,int16,int32,int64,float32,float64,uint,uint8,uint16,uint32,uint64:
		return p == 0
	default:
		return false
	}
}
`

const daoInterface = `
type I{{.Name}}Dao interface {
	Insert({{.PrivateName}} *entity.{{.Name}}) error
	Update({{.PrivateName}} *entity.{{.Name}}) error
	Query({{.PrimaryKey.PrivateName}} {{.PrimaryKey.Tp}}) (*entity.{{.Name}},error)
	Delete({{.PrimaryKey.PrivateName}} {{.PrimaryKey.Tp}}) error
	QueryList(filter entity.{{.Name}}) ([]*entity.{{.Name}},error)
	QueryPage(q entity.{{.Name}}Query) (*entity.{{.Name}}Page,error)
}

var {{.Name}}Dao I{{.Name}}Dao
`

const daoInterfaceBase = `
package dao

import (
	"{{.ProjectName}}/entity"
)

func Init(addr string)  {
	InitDbOrm(addr)

}
`
