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
	ignoreField = "_"
)

func NewDeferCall(expr []ast.Expr) *ast.DeferStmt {
	return &ast.DeferStmt{
		Call: &ast.CallExpr{
			Fun: &ast.CallExpr{
				Fun:  DefaultSelectorExpr,
				Args: expr,
			},
		},
	}
}

func NewCallParams(expr []ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{
		Fun:  DefaultSelectorExpr,
		Args: expr,
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
			deferStmt, exist := checkExistSameDeferCall(node)
			if !exist {
				stmts := node.Body.List
				statList := make([]ast.Stmt, len(stmts)+1)
				copy(statList[1:], stmts)
				statList[0] = NewDeferCall(getFunDeclParams(node))
				node.Body.List = statList
			} else {
				deferStmt.Call.Fun = NewCallParams(getFunDeclParams(node))
			}
			added = true
		}
		return true
	}, nil)
	return fileSet, astFile, added
}

// getFunDeclParams get the func declares params
func getFunDeclParams(node *ast.FuncDecl) (params []ast.Expr) {
	if node.Type == nil || node.Type.Params == nil || len(node.Type.Params.List) == 0 {
		return
	}
	for _, i := range node.Type.Params.List {
		for _, d := range i.Names {
			if d.Name != ignoreField {
				params = append(params, &ast.Ident{Name: d.Name})
			}
		}
	}
	return
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
