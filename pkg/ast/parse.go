package ast

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
)

var (
	DefaultSelIndent = &ast.Ident{
		Name: "Print",
	}
	DefaultXIdent = &ast.Ident{
		Name: "trace_go",
	}
	DefaultSelectorExpr = &ast.SelectorExpr{
		X:   DefaultXIdent,
		Sel: DefaultSelIndent,
	}
	DefaultImportPath = "github.com/samxiaogo/trace_go"
	ignoreField       = "_"
)

func NewSelector(tokenPos token.Pos) *ast.SelectorExpr {
	return &ast.SelectorExpr{X: &ast.Ident{
		NamePos: tokenPos,
		Name:    DefaultXIdent.Name,
	},
		Sel: &ast.Ident{
			NamePos: tokenPos,
			Name:    DefaultSelIndent.Name,
		},
	}
}

func NewDeferCall(expr []ast.Expr, token token.Pos) *ast.DeferStmt {
	return &ast.DeferStmt{
		Defer: token,
		Call: &ast.CallExpr{
			Fun: &ast.CallExpr{
				Fun:    NewSelector(token),
				Args:   expr,
				Lparen: token,
				Rparen: token,
			},
			Lparen: token,
			Rparen: token,
		},
	}
}

func ParseAndAdd(fileName string) []byte {
	fileSet, astFile, added := addDeferMethod(parse(fileName))
	if added {
		astutil.AddImport(fileSet, astFile, DefaultImportPath)
	}
	buf := bytes.NewBufferString("")
	err := format.Node(buf, fileSet, astFile)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func parse(fileName string) (*token.FileSet, *ast.File) {
	fileSet := token.NewFileSet()
	astFile, err := parser.ParseFile(fileSet, fileName, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	return fileSet, astFile
}

func addDeferMethod(fileSet *token.FileSet, astFile *ast.File) (*token.FileSet, *ast.File, bool) {
	var added = false
	astutil.Apply(astFile, func(cursor *astutil.Cursor) bool {
		if node, ok := cursor.Node().(*ast.FuncDecl); ok {
			var (
				deferStmt *ast.DeferStmt
				exist     bool
			)
			params := getFunDeclParams(node, node.Pos())
			deferStmt, exist = checkExistSameDeferCall(node)
			if !exist {
				stmts := node.Body.List
				statList := make([]ast.Stmt, len(stmts)+1)
				copy(statList[1:], stmts)
				statList[0] = NewDeferCall(params, node.Pos())
				node.Body.List = statList
				added = true
			} else if deferStmt != nil {
				if len(params) > 0 {
					deferStmt = NewDeferCall(params, node.Pos())
					added = true
				}
			}
		}
		return true
	}, nil)
	return fileSet, astFile, added
}

// getFunDeclParams get the func declares params
func getFunDeclParams(node *ast.FuncDecl, tokenPos token.Pos) (params []ast.Expr) {
	return
	// todo make it as cmd param
	//if node.Type == nil || node.Type.Params == nil || len(node.Type.Params.List) == 0 {
	//	return
	//}
	//for _, i := range node.Type.Params.List {
	//	for _, d := range i.Names {
	//		if d.Name != ignoreField {
	//			params = append(params, &ast.Ident{Name: d.Name,NamePos: tokenPos})
	//		}
	//	}
	//}
	//return
}

func checkExistSameDeferCall(node *ast.FuncDecl) (*ast.DeferStmt, bool) {
	stmts := node.Body.List
	for _, stmt := range stmts {
		if v, ok := isDeferCallStmt(stmt); ok {
			return v, true
		}
	}
	return nil, false
}

// isDeferCallStmt check is the defer call  todo:check the param is import renamed
func isDeferCallStmt(stmt ast.Stmt) (deferStmt *ast.DeferStmt, exist bool) {
	// is defer stat
	v, ok := stmt.(*ast.DeferStmt)
	if !ok {
		return
	}
	// is defer call()
	c, ok := v.Call.Fun.(*ast.CallExpr)
	if !ok {
		return
	}
	// is selector
	s, ok := c.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	// prefix is the same
	if s.Sel.Name != DefaultSelectorExpr.Sel.Name {
		return
	}
	i, ok := s.X.(*ast.Ident)
	if !ok {
		return
	}
	return v, i.Name == DefaultXIdent.Name
}
