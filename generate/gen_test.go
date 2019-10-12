package generate

import (
	"fmt"
	"github.com/myadamtest/adam/utils"
	"io/ioutil"
	"runtime"
	"testing"
)

func TestGopath(t *testing.T) {
	fmt.Println(runtime.GOROOT())
	fmt.Println(runtime.GOARCH)
	fmt.Println(runtime.GOOS)
	fmt.Println(runtime.GOMAXPROCS(0))
}

func TestParse(t *testing.T) {
	b, _ := ioutil.ReadFile("../grpcservice/greeter_server.go")
	fmt.Println(string(b))
}

func TestCharact(t *testing.T) {
	//str :="abcdefg"
	//
	//for _,gh := range str {
	//	fmt.Println(string(gh))
	//}
	//
	//c := &ParseNodeContext{}
	//
	//obj := reflect.ValueOf(c)
	//m := obj.Method(0)

	//reflect.ValueOf("d")
	//fmt.Println(m.)

	//b,err := ioutil.ReadFile("../protofile/social_relations_service.proto")
	//if err!=nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Println(string(b))

	str := utils.FirstToUpper("agfdre")
	fmt.Println(str)
}

type INode interface {
	Get(out interface{})
	Parse(char int) bool
}

type MethodNode struct {
	MethodName string
	// ... 参数，返回值等设置
}

func (n *MethodNode) Get(out interface{}) {

}

func (n *MethodNode) Parse(char int) bool {
	return false
}

const (
	nodeTypeInit = iota
	nodeTypeMethod
	nodeTypeStr
)

type ParseNodeContext struct {
	Tp        int // 类型
	DoingNode INode
	NodeList  []INode
	//NextNode []*ParseNodeContext // 子节点集合，先做方法一个的就行了。
}

// 怎么切换？
func (pnc *ParseNodeContext) Parse(char int) {
	switch pnc.Tp {
	case nodeTypeMethod:
		_ = pnc.DoingNode.Parse(char)
		/**
		if end add to node list.
		set
		*/
	case nodeTypeInit:
	default:
		fmt.Println("default")
	}
}

var specialCharList = [...]SpecialChar{
	{"\"", "\""},
	{"/*", "*/"},
	{"{", "}"},
	{"(", ")"},
}

type SpecialChar struct {
	Start string
	End   string
}
