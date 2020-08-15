package service

import (
	"github.com/go-xorm/xorm"
	"github.com/myadamtest/adam/gen/service/types"
)

type xxService struct {
	engine xorm.Engine
}

func newXxService() *xxService {
	return &xxService{}
}

func (svc *xxService) Insert(entity *types.Xx) (*types.Xx, error) {
	session := svc.engine
	defer session.Close()

	//todo getID

	_, err := session.InsertOne(entity)
	if err != nil {
		//todo  print err
		return nil, err
	}
	return entity, nil
}

func (svc *xxService) Get(id int64) (*types.Xx, error) {
	session := svc.engine
	defer session.Close()

	entity := new(types.Xx)
	_, err := session.ID(id).Get(entity)
	if err != nil {
		//todo  print err
		return nil, err
	}
	return entity, nil
}

func (svc *xxService) PageQuery(pagination *types.Pagination) ([]*types.Xx, error) {
	session := svc.engine
	defer session.Close()

	var list []*types.Xx

	session.Limit(int(pagination.Count), int(pagination.Page)*int(pagination.Count))
	total, err := session.FindAndCount(&list)
	if err != nil {
		//todo  print err
		return nil, err
	}
	pagination.Total = int32(total)
	return list, err
}

func (svc *xxService) Update(entity *types.Xx) error {
	session := svc.engine
	defer session.Close()

	_, err := session.Where("id = ?", entity.Id).Update(entity)
	if err != nil {
		//todo  print err
		return err
	}

	return nil
}

func (svc *xxService) Delete(id int64) error {
	session := svc.engine
	defer session.Close()

	_, err := session.Id(id).Delete(&types.Xx{})
	if err != nil {
		//todo  print err
		return err
	}
	return nil
}
