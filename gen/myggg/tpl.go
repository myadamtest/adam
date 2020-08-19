package main

import (
	"fmt"
)

type myLog struct {
}

func (log *myLog) Println(arg ...interface{}) {
	fmt.Println(arg)
}

type myEngine struct {
}

func (k *myEngine) TEST() {

}

func (k myEngine) TEST1() {

}

func (e *myEngine) Insert(data interface{}) (int64, error) {
	return 0, nil
}

type xxService struct {
	log    *myLog
	engine *myEngine
	iee    myEngine
}

type A struct {
}

func (x *xxService) Create() error {
	mkdg := A{}
	id, jjf := x.engine.Insert(mkdg)

	if jjf != nil {
		x.log.Println("err:%s,%s", jjf, mkdg)

		mkl := func(i int64) int64 {
			return i + 1
		}
		mkl(1)

		return jjf
	}

	fmt.Println(id)
	return nil
}

func Mk() {

}
