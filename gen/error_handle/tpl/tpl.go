package foo

import "fmt"

type myLog struct {
}

func (log *myLog) Println(arg ...interface{}) {
	fmt.Println(arg)
}

type myEngine struct {
}

func (e *myEngine) Insert(data interface{}) (int64, error) {
	return 0, nil
}

type xxService struct {
	log    *myLog
	engine *myEngine
}

type A struct {
}

func (x *xxService) Create() error {
	mkdg := A{}
	id, jjf := x.engine.Insert(mkdg)
	if jjf != nil {
		x.log.Println("err:%s,%s", jjf, mkdg)
		return jjf
	}

	fmt.Println(id)
	return nil
}
