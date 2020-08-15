package main

import (
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/myadamtest/adam/dbinfo"
	"github.com/myadamtest/adam/gen/service"
)

func init() {
	flag.Parse()
}

func main() {
	fp := dbinfo.ToGogoProto(dbinfo.GetTables())
	//gogoproto.Gen(fp)
	service.Gen(fp)
}
