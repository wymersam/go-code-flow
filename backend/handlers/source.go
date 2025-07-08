package handlers

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
)

func getFuncSource(node *ast.FuncDecl, fset *token.FileSet) (string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, node)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
