package main

import (
	"fmt"
	"github.com/pkg/errors"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ast/astutil"
	"io/ioutil"
	"math"
	"reflect"
)

func main() {
	src, _ := ioutil.ReadFile("./gen/myggg/tpl.go")

	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, "", string(src), 0)
	if err != nil {
		fmt.Println("Can't parse file", err)
	}

	info := &types.Info{Types: make(map[ast.Expr]types.TypeAndValue)}
	conf := types.Config{Importer: importer.Default()}
	if _, err = conf.Check("", fs, []*ast.File{file}, info); err != nil {
		fmt.Println(err)
		return
	}

	global := 0
	logCtx := &LogContext{}
	astutil.Apply(file, func(cursor *astutil.Cursor) bool {
		logCtx.index++
		global++
		if 310 == global {
			fmt.Println("")
		}
		fmt.Println(logCtx.index, global, cursor.Node())
		switch n := cursor.Node().(type) {
		case *ast.FuncDecl:
			express, err := GetFuncLogExpress(n)
			if err != nil {

			} else {
				logCtx.funIndex = logCtx.index
				logCtx.PreExpress = express
			}
			return true
		}

		// 表达式为空，证明不可以打印日志
		if logCtx.PreExpress == nil {
			return true
		}

		if t, ok := cursor.Node().(*ast.AssignStmt); !ok {
			return true
		} else if len(t.Rhs) == 0 {
			return true
			// t.Rhs[0] 确定顶只有T0是call?
		} else if callExpr, ok := t.Rhs[0].(*ast.CallExpr); !ok {
			return true
		} else {
			var results *types.Tuple
			if funType, ok := info.Types[callExpr]; !ok {
				return true
			} else if funRealType, ok := funType.Type.(*types.Signature); !ok {
				return true
			} else {
				results = funRealType.Results()
			}

			if results.Len() == 0 {
				return true
			}

			last := results.At(results.Len() - 1)
			if named, ok := last.Type().(*types.Named); !ok {
				return true
			} else if named.Obj() == nil {
				return true
			} else if named.Obj().Id() != "_.error" {
				return true
			}
			// 是错误，打印，只判断最后一个参数是否错误类型

			// 返回参数没有接收
			if len(t.Lhs) == 0 {
				return true
			}

			ev, ok := t.Lhs[results.Len()-1].(*ast.Ident)
			if !ok {
				return true
			}
			logAssign := GetLogPrintStmt(callExpr.Args, logCtx, ev.String())
			logCtx.CurrentCursor = cursor
			logCtx.LogAssign = logAssign
		}
		return true
	}, func(cursor *astutil.Cursor) bool {
		logCtx.index--
		global++
		//fmt.Println(logCtx.index,global,cursor.Node())
		// 一个符合条件的函数结束
		if logCtx.index == logCtx.funIndex {
			fmt.Println(logCtx.index)
			logCtx.PreExpress = nil
			logCtx.funIndex = math.MaxInt32
		}
		return true
	})
}

func GetLogPrintStmt(args []ast.Expr, ctx *LogContext, errValuableName string) *ast.ExprStmt {
	newArgs := make([]ast.Expr, 2, len(args)+2)
	for _, arg := range args {
		a := arg.(*ast.Ident)
		aaa := &ast.Ident{Name: a.Name}
		newArgs = append(newArgs, aaa)
	}
	newArgs[0] = &ast.BasicLit{
		Kind:  token.STRING,
		Value: "\"err:%s,%s\"",
	}

	newArgs[1] = &ast.Ident{Name: errValuableName}

	a1 := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X: &ast.SelectorExpr{
					X: &ast.Ident{Name: ctx.PreExpress.Express},
				},
				Sel: &ast.Ident{Name: "Println"},
			},
			Args: newArgs,
		},
	}
	return a1
}

func GetFuncLogExpress(funNode *ast.FuncDecl) (*LogExpress, error) {
	if funNode == nil || funNode.Recv == nil {
		return nil, errors.New("recv is empty")
	}
	if len(funNode.Recv.List) == 0 {
		return nil, errors.New("recv list is empty")
	}

	recv := funNode.Recv.List[0]
	if len(recv.Names) == 0 {
		return nil, errors.New("recv name is empty")
	}

	var ident *ast.Ident
	if starExpr, ok := recv.Type.(*ast.StarExpr); !ok {
		if i, ok := recv.Type.(*ast.Ident); !ok {
			return nil, errors.New("this fun not have correct recv")
		} else {
			ident = i
		}
	} else if x, ok := starExpr.X.(*ast.Ident); !ok {
		return nil, errors.New("this fun not have correct recv")
	} else {
		ident = x
	}

	if ident.Obj == nil || ident.Obj.Decl == nil {
		return nil, errors.New("this fun not have correct recv")
	}

	var fieldList []*ast.Field
	if ts, ok := ident.Obj.Decl.(*ast.TypeSpec); !ok {
		return nil, errors.New("this fun not have correct TypeSpec")
	} else if st, ok := ts.Type.(*ast.StructType); !ok {
		return nil, errors.New("this fun not have correct StructType")
	} else if st.Fields == nil {
		return nil, errors.New("this fun not have correct StructType.Fields")
	} else {
		fieldList = st.Fields.List
	}
	if len(fieldList) == 0 {
		return nil, errors.New("this struct not have field")
	}

	// 定制
	logPackage := "myLog"
	fullName := ""
	for _, field := range fieldList {
		if se, ok := field.Type.(*ast.StarExpr); !ok {
			if i, ok := field.Type.(*ast.Ident); !ok {
				if selector, ok := field.Type.(*ast.SelectorExpr); ok {
					fullName = SelectorExpressToString(selector)
				}
			} else {
				fullName = i.Name
			}
		} else if x, ok := se.X.(*ast.Ident); !ok {
			if selector, ok := se.X.(*ast.SelectorExpr); ok {
				fullName = SelectorExpressToString(selector)
			}
		} else {
			fullName = x.Name
		}

		if fullName == logPackage {
			break
		}
	}
	if fullName != logPackage {
		return nil, errors.New("没有导入错误包")
	}

	return &LogExpress{Express: fmt.Sprintf("%s.%s", recv.Names[0], fullName)}, nil
}

func SelectorExpressToString(selector *ast.SelectorExpr) string {
	result := ""
	if s, ok := selector.X.(*ast.SelectorExpr); ok {
		result = SelectorExpressToString(s)
	} else if i, ok := selector.X.(*ast.Ident); ok {
		result = i.Name
	}

	if result == "" {
		return selector.Sel.Name
	}
	return fmt.Sprintf("%s.%s", result, selector.Sel.Name)
}

type LogExpress struct {
	Express string
}

type LogContext struct {
	PreExpress    *LogExpress // 日志前缀表达式
	CurrentCursor *astutil.Cursor
	LogAssign     *ast.ExprStmt
	index         int64
	funIndex      int64
}

func printType(i interface{}) {
	fmt.Println(reflect.TypeOf(i), ">>>>>>>>>>>>>>>")
}
