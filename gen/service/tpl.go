package service

const serviceTpl = `
package service

import (
	"github.com/go-xorm/xorm"
	"github.com/myadamtest/adam/gen/service/types"
)

type {{firstLowerName .Name}}Service struct {
	engine xorm.Engine
}

func new{{firstUpName .Name}}Service() *{{firstLowerName .Name}}Service {
	return &{{firstLowerName .Name}}Service{}
}

func (svc *{{firstLowerName .Name}}Service) Insert(entity *types.{{firstUpName .Name}}) (*types.{{firstUpName .Name}},error){
	session := svc.engine
	defer session.Close()

	//todo getID

	_,err := session.InsertOne(entity)
	if err!= nil {
		//todo  print err
		return nil,err
	}
	return entity,nil
}

func (svc *{{firstLowerName .Name}}Service) Get(id int64) (*types.{{firstUpName .Name}},error)  {
	session := svc.engine
	defer session.Close()

	entity := new(types.{{firstUpName .Name}})
	_,err := session.ID(id).Get(entity)
	if err!= nil {
		//todo  print err
		return nil,err
	}
	return entity,nil
}

func (svc *{{firstLowerName .Name}}Service) PageQuery(pagination *types.Pagination) ([]*types.{{firstUpName .Name}},error) {
	session := svc.engine
	defer session.Close()

	var list []*types.{{firstUpName .Name}}

	session.Limit(int(pagination.Count),int(pagination.Page) * int(pagination.Count))
	total,err := session.FindAndCount(&list)
	if err!= nil {
		//todo  print err
		return nil,err
	}
	pagination.Total = int32(total)
	return list,err
}

func (svc *{{firstLowerName .Name}}Service) Update(entity *types.{{firstUpName .Name}}) error  {
	session := svc.engine
	defer session.Close()

	_,err := session.Where("id = ?",entity.Id).Update(entity)
	if err!= nil {
		//todo  print err
		return err
	}
	
	return nil
}

func (svc *{{firstLowerName .Name}}Service) Delete(id int64) error {
	session := svc.engine
	defer session.Close()

	_,err := session.Id(id).Delete(&types.{{firstUpName .Name}}{})
	if err!= nil {
		//todo  print err
		return err
	}
	return nil
}
`
