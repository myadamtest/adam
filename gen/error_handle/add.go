package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ast/astutil"
	"io/ioutil"
	"os"
)

func main() {
	isError := func(v ast.Expr, info *types.Info) bool {
		if n, ok := info.TypeOf(v).(*types.Named); ok {
			o := n.Obj()
			return o != nil && o.Pkg() == nil && o.Name() == "error"
		}
		return false
	}

	src, _ := ioutil.ReadFile("./gen/error_handle/tpl/tpl.go")

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

	astutil.Apply(file, func(cursor *astutil.Cursor) bool {
		switch n := cursor.Node().(type) {
		case *ast.FuncDecl:
			recv := n.Recv.List[0]

			stFields := recv.Type.(*ast.StarExpr).X.(*ast.Ident).Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType).Fields.List
			if len(stFields) == 0 {
				return true
			}

			logVariableName := ""

			for _, f := range stFields {
				if f.Type.(*ast.StarExpr).X.(*ast.Ident).Name == "myLog" {
					logVariableName = f.Names[0].Name
				}
			}

			for index, line := range n.Body.List {
				switch n2 := line.(type) {
				case *ast.AssignStmt:
					i, ok := n2.Lhs[0].(*ast.Ident)
					if ok {
						if i.Obj != nil {
							asss := i.Obj.Decl.(*ast.AssignStmt)

							ce, ok := asss.Rhs[0].(*ast.CallExpr)
							if !ok {
								break
							}

							isErr := isError(ce, info)
							fmt.Println("OOOOOOOOOOO+", isErr, ce.Args)
							args := ce.Args
							se := ce.Fun.(*ast.SelectorExpr)

							results := info.Types[se].Type.(*types.Signature).Results()

							// len > 0
							t := results.At(results.Len() - 1)
							isE := t.Type().(*types.Named).Obj().Id() == "_.error"

							//errorType := t.Type().(*types.Named).Obj().Type()
							fmt.Println(isE, asss.Lhs[1])

							//obj := info.Uses[se.Sel]
							//fmt.Println(info.Types[se].Type.(*types.Signature).Results().At(1).Type().String(),"||||||||||",info)

							fmt.Println(se.X)

							// 后面还有行 而其是
							if index < len(n.Body.List)-1 {
								if ifst, ok := n.Body.List[index+1].(*ast.IfStmt); ok {
									be := ifst.Cond.(*ast.BinaryExpr)
									if be.X.(*ast.Ident).String() == asss.Lhs[1].(*ast.Ident).String() && be.Op.String() == "!=" && be.Y.(*ast.Ident).String() == "nil" {
										// 后面有if并且是判断if的
										existLog := false
										for _, ib := range ifst.Body.List {
											if esTemp, ok := ib.(*ast.ExprStmt); ok {
												if ce, ok := esTemp.X.(*ast.CallExpr); ok {
													if fe, ok := ce.Fun.(*ast.SelectorExpr); ok {
														isXLog := false
														if se, ok := fe.X.(*ast.SelectorExpr); ok {
															if se.X.(*ast.Ident).String() == "x" && se.Sel.String() == "log" {
																isXLog = true
															}
														}
														if isXLog && fe.Sel.String() == "Println" {
															existLog = true
														}
													}
												}
											}
										}

										if existLog {
											continue
										}

										//newArgs := make([]ast.Expr,0)
										newArgs := make([]ast.Expr, 2, len(args)+2)
										for _, arg := range args {
											a := arg.(*ast.Ident)
											fmt.Println(a.Name, info.TypeOf(a).String())
											aaa := &ast.Ident{Name: a.Name}
											newArgs = append(newArgs, aaa)
										}
										newArgs[0] = &ast.BasicLit{
											Kind:  token.STRING,
											Value: "\"err:%s,%s\"",
										}

										newArgs[1] = &ast.Ident{Name: asss.Lhs[1].(*ast.Ident).String()}

										// add log
										fmt.Println(n.Recv.List[0], "<<<<<<<<")
										a1 := ast.ExprStmt{
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.SelectorExpr{
														X:   &ast.Ident{Name: "x"},
														Sel: &ast.Ident{Name: logVariableName},
													},
													Sel: &ast.Ident{Name: "Println"},
												},
												Args: newArgs,
											},
										}

										ifst.Body.List = append([]ast.Stmt{&a1}, ifst.Body.List...)
									}
								}
							}
						}
					}

				case *ast.IfStmt:
					//fmt.Println(n2.Pos())
					//be := n2.Cond.(*ast.BinaryExpr)
					//fmt.Println(be.X,be.Op.String(),be.Y)
				}

			}
		}

		return true
	}, func(cursor *astutil.Cursor) bool {
		return true
	})

	printer.Fprint(os.Stdout, fs, file)
}
