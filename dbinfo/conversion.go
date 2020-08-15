package dbinfo

import (
	"github.com/myadamtest/adam/gen/gogoproto"
	"github.com/myadamtest/adam/gen/proto"
	"strings"
	"xorm.io/core"
)

// sql类型转换proto类型
var typeMap = map[string]string{
	"int":       "int64",
	"tinyint":   "int32",
	"smallint":  "int32",
	"decimal":   "float64",
	"varchar":   "string",
	"char":      "string",
	"timestamp": "string",
}

func ToGogoProto(tables []*core.Table) *gogoproto.ProtoFile {
	pf := &gogoproto.ProtoFile{}
	// 定制化参数
	pf.Syntax = "proto3"
	pf.PackageName = "types"

	for _, t := range tables {
		message := proto.MessageElement{
			Name: t.Name, // tp camel form
		}

		for _, c := range t.Columns() {

			typ, _ := proto.NewScalarDataType(typeMap[strings.ToLower(c.SQLType.Name)])
			field := proto.FieldElement{
				Name: c.Name, // tp camel form
				Type: typ,
			}
			message.Fields = append(message.Fields, field)
		}

		pf.Messages = append(pf.Messages, message)
	}

	return pf
}
