package dbinfo

import (
	"flag"
	"fmt"
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

var dbAddr = flag.String("db-addr", "", "get db info")
var tables []*core.Table

func GetTables() []*core.Table {
	if tables != nil {
		return tables[:]
	}

	err := initTables(*dbAddr)
	fmt.Println(err)

	return tables[:]
}

func initTables(addr string) error {
	engine, err := xorm.NewEngine("mysql", addr)
	if err != nil {
		return err
	}

	tables, err = engine.DBMetas()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
