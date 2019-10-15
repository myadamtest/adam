package dbgenerate

import (
	"fmt"
	"github.com/myadamtest/adam/generate"
)

var generateFuns = []func(*structInfo) error{
	generateStruct, generateDao, generateIDao, generateService, generateIService, generaGrpc, generateGrpcWithImpl,
}

func GenCode(addr string) error {
	//创建数据库
	sc, err := newSqlCli(addr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//获取数据库表信息
	tables, err := sc.getTables()
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = generateCommonOrm()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, t := range tables {
		if t == nil {
			continue
		}
		si := tableConversion2Struct(t)

		for i, f := range generateFuns {
			err = f(si)
			if err != nil {
				fmt.Println(err, i)
			}
		}
	}

	generate.GrpcGenerate()

	return nil
}
