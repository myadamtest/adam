package generate

import (
	"bufio"
	"github.com/bookrun-go/fileutils/fileinfo"
	"github.com/myadamtest/adam/utils"
	"io"
	"os"
	"regexp"
	"strings"
)

type RpcParam struct {
	Name string
	Ty   int // 0普通参数，流类型
}

type RpcServiceInterfaceInfo struct {
	Name          string
	RequestParam  RpcParam
	ResponseParam RpcParam
}

type RpcServiceInfo struct {
	Name        string
	FileName    string
	ProjectName string
	PackageName string
	Interfaces  []RpcServiceInterfaceInfo
}

func GetProtoStruct(fp string) (*RpcServiceInfo, error) {
	fi, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	projectName, err := utils.GetProjectName()
	if err != nil {
		return nil, err
	}
	simpleName := fileinfo.GetFileSimpleName(fp)

	br := bufio.NewReader(fi)
	rpcServiceInfo := RpcServiceInfo{}
	rpcServiceInfo.Interfaces = make([]RpcServiceInterfaceInfo, 0)
	rpcServiceInfo.FileName = simpleName

	rpcServiceInfo.ProjectName = projectName

	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

		lineStr := string(a)

		newLine := strings.Trim(lineStr, " ")

		if strings.HasPrefix(newLine, "package ") {
			newLine = deleteExtraSpace(newLine)
			newLine = newLine[8:]

			result := strings.Index(newLine, ";")
			rpcServiceInfo.PackageName = strings.Trim(newLine[0:result], " ")
			continue
		}

		if strings.HasPrefix(newLine, "service ") {
			newLine = deleteExtraSpace(newLine)
			newLine = newLine[8:]
			//fmt.Println(newLine)

			result := strings.Index(newLine, " ")
			//fmt.Println(result)
			//
			//fmt.Println("service name =",newLine[0:result])
			rpcServiceInfo.Name = newLine[0:result]
			continue
		}

		if strings.HasPrefix(newLine, "rpc ") {
			newLine = deleteExtraSpace(newLine)
			newLine = newLine[4:]
			//fmt.Println(newLine)

			iInfo := RpcServiceInterfaceInfo{}

			result := strings.Index(newLine, " ")
			//fmt.Println(result)

			iInfo.Name = newLine[0:result]
			iInfo.Name = utils.FirstToUpper(iInfo.Name)
			//fmt.Println("method name =",newLine[0:result])

			newLine = strings.Trim(newLine[result:], " ")

			result = strings.Index(newLine, ")")
			reqStr := strings.Trim(newLine[1:result], " ")

			reqParam := RpcParam{}
			if strings.HasPrefix(reqStr, "stream ") {
				reqParam.Name = strings.Trim(strings.TrimLeft(reqStr, "stream "), " ")
				reqParam.Ty = 1
			} else {
				reqParam.Name = strings.Trim(reqStr, " ")
			}
			reqParam.Name = utils.FirstToUpper(reqParam.Name)

			iInfo.RequestParam = reqParam

			newLine = newLine[result+11:]
			result = strings.Index(newLine, ")")
			respStr := strings.Trim(newLine[0:result], " ")

			respParam := RpcParam{}
			if strings.HasPrefix(respStr, "stream ") {
				respParam.Name = strings.Trim(strings.TrimLeft(respStr, "stream "), " ")
				respParam.Ty = 1
			} else {
				respParam.Name = strings.Trim(respStr, " ")
			}
			respParam.Name = utils.FirstToUpper(respParam.Name)
			iInfo.ResponseParam = respParam

			rpcServiceInfo.Interfaces = append(rpcServiceInfo.Interfaces, iInfo)
			continue
		}
	}
	return &rpcServiceInfo, nil
}

func deleteExtraSpace(s string) string {
	//删除字符串中的多余空格，有多个空格时，仅保留一个空格
	s1 := strings.Replace(s, "  ", " ", -1)      //替换tab为空格
	regstr := "\\s{2,}"                          //两个及两个以上空格的正则表达式
	reg, _ := regexp.Compile(regstr)             //编译正则表达式
	s2 := make([]byte, len(s1))                  //定义字符数组切片
	copy(s2, s1)                                 //将字符串复制到切片
	spc_index := reg.FindStringIndex(string(s2)) //在字符串中搜索
	for len(spc_index) > 0 {                     //找到适配项
		s2 = append(s2[:spc_index[0]+1], s2[spc_index[1]:]...) //删除多余空格
		spc_index = reg.FindStringIndex(string(s2))            //继续在字符串中搜索
	}
	return string(s2)
}
