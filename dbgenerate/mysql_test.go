package dbgenerate

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"testing"
	"xorm.io/core"
)

//"user:password@/dbname?charset=utf8&parseTime=True&loc=Local"
func TestGetTableInfo(t *testing.T) {
	engine, err := xorm.NewEngine("mysql", "root:1234@tcp(39.107.244.239:3316)/ddc?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
		return
	}

	tables, err := engine.DBMetas()
	if err != nil {
		fmt.Println(err)
		return
	}

	f := func(t *core.Table) {
		fmt.Println(t.Name)
		fmt.Println(t.Columns()[0])
	}

	f(tables[0])
}
