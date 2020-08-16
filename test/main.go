package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ast/astutil"
	"reflect"
	"strings"
)

const templateGo = `
package foo

import (
  "net/http"
  "fmt"
)

var y6 int

type Foo struct {
}

func (f *Foo) Method(req *http.Request) *http.Response {
	_ = new(Foo)
	r := f.ttMethod222(req,1)
	fmt.Println(r)
	return nil
}

func (f *Foo) ttMethod222(req *http.Request, hu int) *http.Response {
	_ = new(Foo)
	return nil
}

func plus(i int) {
	
}

`

func main() {
	isError := func(v ast.Expr, info *types.Info) bool {
		if n, ok := info.TypeOf(v).(*types.Named); ok {
			o := n.Obj()
			return o != nil && o.Pkg() == nil && o.Name() == "error"
		}
		return false
	}

	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, "", templateGo, 0)
	if err != nil {
		fmt.Println("Can't parse file", err)
	}

	// We extract type info
	info := &types.Info{Types: make(map[ast.Expr]types.TypeAndValue)}
	conf := types.Config{Importer: importer.Default()}
	if _, err = conf.Check("", fs, []*ast.File{file}, info); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(isError)

	astutil.Apply(file, func(cursor *astutil.Cursor) bool {
		switch n := cursor.Node().(type) {
		case *ast.FuncDecl:
			for _, line := range n.Body.List {
				switch n2 := line.(type) {
				case *ast.AssignStmt:
					i, ok := n2.Lhs[0].(*ast.Ident)
					if ok {
						fmt.Println(">>>>", i.Obj)
						if i.Obj != nil {
							asss := i.Obj.Decl.(*ast.AssignStmt)

							ce := asss.Rhs[0].(*ast.CallExpr)
							isErr := isError(ce, info)
							fmt.Println("OOOOOOOOOOO+", isErr, ce)
							se := ce.Fun.(*ast.SelectorExpr)

							//obj := info.Uses[se.Sel]
							fmt.Println(info.Types[se].Type.(*types.Signature).Results().At(0).Type().String(), "||||||||||", info)

							fmt.Println(se.X)
						}
					}

				}

			}
		}

		return true
	}, func(cursor *astutil.Cursor) bool {
		return true
	})

	//file.Name.Name = "bar111"
	//
	//r := &Renamer{"Foo", "Bar"}
	//ast.Walk(r, file)
	//
	//printer.Fprint(os.Stdout, fs, file)
}

type Renamer struct {
	find    string
	replace string
}

func (r *Renamer) Visit(node ast.Node) (w ast.Visitor) {
	if node != nil {
		switch n := node.(type) {
		case *ast.FuncDecl:
			if n.Recv != nil && n.Recv.List != nil && len(n.Recv.List) > 0 {
				fmt.Println(len(n.Recv.List), ">>>>>>>>>>>>>>>>")
				field := n.Recv.List[0]
				typ := field.Type.(*ast.StarExpr).X.(*ast.Ident).Name
				if typ == r.find {
					field.Names[0].Name = strings.ToLower(r.replace[0:1])
				}
			}
		case *ast.Ident:
			if n.Name == r.find {
				n.Name = r.replace
			}
		case *ast.CaseClause:

		default:
			fmt.Println(reflect.TypeOf(n).String())
		}
	}
	return r
}
