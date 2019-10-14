package dbgenerate

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func generateDao(si *structInfo) error {
	err := os.Mkdir("./dao", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	daoFilename := fmt.Sprintf("./dao/%s.go", si.TableName)

	err = executeTemplate(daoTemplate, daoFilename, si, true)
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
	_, err := os.Stat(filename)

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
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
	contentArr := strings.Split(content, "InitDbOrm(addr)")

	if len(contentArr) != 2 {
		return errors.New("unknown err")
	}

	content = contentArr[0] + "InitDbOrm(addr)\n\t" + info.Name + "Dao = new" + info.Name + "Dao()" + contentArr[1]

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

func executeTemplate(tpl, fullFilename string, info *structInfo, append bool) error {
	return executeTemplateWithFuncs(tpl, fullFilename, info, append, nil)
}

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

const daoTemplate = `
package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/myadamtest/gobase/logkit"
	"{{.ProjectName}}/grpcservice/pb/pb"
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
	err := dao.get{{.Name}}Db().Update({{.PrivateName}}).Error
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
		{{.Name}}:pb.{{.Name}}{
			{{.PrimaryKey.Name}}:{{.PrimaryKey.PrivateName}},
		},
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
	err := dao.get{{.Name}}Db().Where(filter).Find(result).Error
	if err!=nil {
		logkit.Errorf("query list {{.PrivateName}} err:%s",err)
		return nil,err
	}
	return result, nil
}

func (dao *{{.PrivateName}}Dao) QueryPage(q pb.{{.Name}}PageRequest) (* pb.{{.Name}}PageResponse,error)  {
	page := &pb.{{.Name}}PageResponse{}
	page.Page = q.Page

	db := dao.get{{.Name}}Db().Where(q.{{.Name}})
	err := db.Count(&page.Page.Total).Error
	if err!= nil {
		logkit.Errorf("get {{.PrivateName}} total err:%s",err)
		return nil,err
	}

	err = db.Offset((page.Page.PageNo -1) * page.Page.PageSize).Limit( page.Page.PageSize).Find(& page.{{.Name}}List).Error
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
`

const daoInterface = `
type I{{.Name}}Dao interface {
	Insert({{.PrivateName}} *entity.{{.Name}}) error
	Update({{.PrivateName}} *entity.{{.Name}}) error
	Query({{.PrimaryKey.PrivateName}} {{.PrimaryKey.Tp}}) (*entity.{{.Name}},error)
	Delete({{.PrimaryKey.PrivateName}} {{.PrimaryKey.Tp}}) error
	QueryList(filter entity.{{.Name}}) ([]*entity.{{.Name}},error)
	QueryPage(q pb.{{.Name}}PageRequest) (*pb.{{.Name}}PageResponse,error)
}

var {{.Name}}Dao I{{.Name}}Dao
`

const daoInterfaceBase = `
package dao

import (
	"{{.ProjectName}}/entity"
	"{{.ProjectName}}/grpcservice/pb/pb"
)

func Init(addr string)  {
	InitDbOrm(addr)

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
