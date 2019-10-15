package dbgenerate

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	serviceInitSplit  = "func Init() {"
	serviceContentFmt = "%s%s\n\t%sService = new%sService()%s"
)

func generateService(info *structInfo) error {
	err := os.Mkdir("./service", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	serviceFilename := fmt.Sprintf("./service/%s.go", info.TableName)

	err = executeTemplate(serviceTemplate, serviceFilename, info, true)
	if err != nil {
		return err
	}

	return nil
}

func generateIService(info *structInfo) error {
	filename := "./service/IService.go"

	exist, err := serviceIsExist(filename)
	if err != nil {
		return err
	}

	if !exist {
		err = executeTemplate(serviceInterfaceBase, filename, info, false)
		if err != nil {
			return err
		}
	}

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	content := string(b)
	contentArr := strings.Split(content, serviceInitSplit)

	if len(contentArr) != 2 {
		return errors.New("unknown err")
	}

	content = fmt.Sprintf(serviceContentFmt, contentArr[0], serviceInitSplit, info.Name, info.Name, contentArr[1])
	err = ioutil.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return err
	}

	err = executeTemplate(serviceInterface, filename, info, true)
	if err != nil {
		return err
	}

	return nil
}

func serviceIsExist(filename string) (bool, error) {
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
	if str == "package service" {
		return false, nil
	}
	return true, nil
}

const serviceTemplate = `
package service

import (
	"{{.ProjectName}}/dao"
	"{{.ProjectName}}/entity"
)

type {{.PrivateName}}Service struct {

}

func new{{.Name}}Service() I{{.Name}}Service {
	return &{{.PrivateName}}Service{}
}

func (serv *{{.PrivateName}}Service) Insert({{.PrivateName}} *entity.{{.Name}}) error  {
	return dao.{{.Name}}Dao.Insert({{.PrivateName}})
}

func (serv *{{.PrivateName}}Service) Update({{.PrivateName}} *entity.{{.Name}}) error  {
	return dao.{{.Name}}Dao.Update({{.PrivateName}})
}

func (serv *{{.PrivateName}}Service) Query({{.PrimaryKey.PrivateName}} {{.PrimaryKey.Tp}}) (*entity.{{.Name}},error)  {
	return dao.{{.Name}}Dao.Query({{.PrimaryKey.PrivateName}})
}

func (serv *{{.PrivateName}}Service) Delete({{.PrimaryKey.PrivateName}} {{.PrimaryKey.Tp}}) error {
	return dao.{{.Name}}Dao.Delete({{.PrimaryKey.PrivateName}})
}

func (serv *{{.PrivateName}}Service) QueryList(filter entity.{{.Name}}) ([]*entity.{{.Name}},error) {
	return dao.{{.Name}}Dao.QueryList(filter)
}

func (serv *{{.PrivateName}}Service) QueryPage(q entity.{{.Name}}Query) (*entity.{{.Name}}Page,error)  {
	return dao.{{.Name}}Dao.QueryPage(q)
}
`

const serviceInterfaceBase = `
package service

import (
	"{{.ProjectName}}/entity"
)

func Init() {

}
`

const serviceInterface = `
type I{{.Name}}Service interface {
	Insert({{.PrivateName}} *entity.{{.Name}}) error
	Update({{.PrivateName}} *entity.{{.Name}}) error
	Query({{.PrimaryKey.PrivateName}} {{.PrimaryKey.Tp}}) (*entity.{{.Name}},error)
	Delete({{.PrimaryKey.PrivateName}} {{.PrimaryKey.Tp}}) error
	QueryList(filter entity.{{.Name}}) ([]*entity.{{.Name}},error)
	QueryPage(q entity.{{.Name}}Query) (*entity.{{.Name}}Page,error)
}

var {{.Name}}Service I{{.Name}}Service
`
