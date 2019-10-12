package mygen

import "fmt"

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

		err = generateStruct(si)
		if err != nil {
			fmt.Println(err)
		}

		err = generateDao(si)
		if err != nil {
			fmt.Println(err)
		}

		err = generateIDao(si)
		if err != nil {
			fmt.Println(err)
		}

		err = generateService(si)
		if err != nil {
			fmt.Println(err)
		}

		err = generateIService(si)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
