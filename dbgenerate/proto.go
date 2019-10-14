package dbgenerate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

func generaGrpc(info *structInfo) error {
	info = copyStructInfo(info)

	err := os.Mkdir("./protofile", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	err = createCommonProto()
	if err != nil {
		return err
	}

	return executeTemplateWithFuncs(protoTemplate, fmt.Sprintf("./protofile/%s.proto", info.TableName), info, false, template.FuncMap{"indexAdd": indexAdd})
}

func indexAdd(i int) int {
	return i + 1
}

func copyStructInfo(info *structInfo) *structInfo {
	nInfo := &structInfo{}
	*nInfo = *info

	str, _ := json.Marshal(info)
	_ = json.Unmarshal(str, nInfo)

	if nInfo.PrimaryKey != nil {
		nInfo.PrimaryKey.Tp = type2Proto(nInfo.PrimaryKey.Tp)
	}

	for _, v := range nInfo.FieldInfos {
		if v == nil {
			continue
		}

		v.Tp = type2Proto(v.Tp)
	}

	return nInfo
}

func type2Proto(tp string) string {
	switch tp {
	case "int", "int8", "int16":
		return "int32"
	case "float32", "float64":
		return "float"
	default:
		return tp
	}
}

func createCommonProto() error {
	filename := "./protofile/common.proto"
	_, err := os.Stat(filename)
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return nil
	}

	return ioutil.WriteFile(filename, []byte(commonProtoTemplate), 0644)
}

const protoTemplate = `
syntax = "proto3";

import "common.proto";

package pb;

service {{.Name}}Service {
  // 增加数据
  rpc Insert ({{.Name}}) returns ({{.Name}});
  // 修改数据
  rpc Update ({{.Name}}) returns ({{.Name}});
  // 根据主键查询
  rpc Query ({{.Name}}PkParamRequest) returns ({{.Name}});
  // 根据主键删除
  rpc Delete ({{.Name}}PkParamRequest) returns ({{.Name}});
  // 根据条件查询
  rpc QueryList ({{.Name}}) returns ({{.Name}}Array);
  // 分页查询
  rpc QueryPage ({{.Name}}PageRequest) returns ({{.Name}}PageResponse);
}

message {{.Name}}PkParamRequest {
  	{{if .PrimaryKey}} {{.PrimaryKey.Tp}} {{.PrimaryKey.Name}} = 1; {{.PrimaryKey.Comment}}
	{{end}}
}

message {{.Name}} {
	{{range $i,$v :=.FieldInfos}} {{$v.Tp}} {{$v.Name}} = {{indexAdd $i}}; {{$v.Comment}}
	{{end}}
}

message {{.Name}}Array {
  repeated {{.Name}} {{.Name}}List = 1;
}

message {{.Name}}PageRequest {
  Page Page = 1;
  {{.Name}} {{.Name}} = 2;
}

message {{.Name}}PageResponse {
  Page Page = 1;
  repeated {{.Name}} {{.Name}}List = 2;
}
`

const commonProtoTemplate = `
syntax = "proto3";

package pb;

message Page {
  int32 PageNo = 1;
  int32 PageSize = 2;
  int64 Total = 3;
}
`
