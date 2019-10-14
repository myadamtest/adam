package main

import (
	"fmt"
	"github.com/myadamtest/adam/dbgenerate"
	"os"
	"testing"
)

func TestCreateStruct(t *testing.T) {
	rm()

	err := dbgenerate.GenCode("dd:2222ddD_@tcp(47.94.168.30:3316)/ddc")
	fmt.Println(err)
}

func TestRemove(t *testing.T) {
	rm()
}

//func TestTable(t *testing.T)  {
//	dao.Init("dd:2222ddD_@tcp(47.94.168.30:3316)/ddc")
//	service.Init()
//
//	myArt := &entity.Article{pb.Article{Title:"dddd"}}
//	err := service.ArticleService.Insert(myArt)
//
//	art,err := service.ArticleService.Query(myArt.Id)
//	if err!=nil {
//		panic(err)
//	}
//	fmt.Println(art.Title,art.Id)
//}

func rm() {
	_ = os.RemoveAll("./dao")
	_ = os.RemoveAll("./service")
	_ = os.RemoveAll("./entity")
	_ = os.RemoveAll("./grpcservice")
	_ = os.RemoveAll("./protofile")

}
