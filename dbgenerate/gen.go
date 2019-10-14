package dbgenerate

import (
	"fmt"
	"github.com/myadamtest/adam/generate"
)

var generateFuns = []func(*structInfo) error{
	generateStruct, generateDao, generateIDao, generateService, generateIService, generaGrpc,
}

func GenCode(addr string) error {
	//创建数据库
	sc, err := newSqlCli(addr)
	if err != nil {
		return err
	}

	//获取数据库表信息
	tables, err := sc.getTables()
	if err != nil {
		return err
	}

	err = generateCommonOrm()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, t := range tables {
		if t == nil {
			fmt.Println(t)
			continue
		}
		si := tableConversion2Struct(t)

		for _, f := range generateFuns {
			err = f(si)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	generate.GrpcGenerate()

	return nil
}
