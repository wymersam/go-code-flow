package handlers

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"

	"github.com/wymersam/goflow/api"
)

// CallGraph represents a function call graph: function => list of called functions
type CallGraph map[string][]string
type FunctionInfo struct {
	Name       string
	Calls      []string
	SourceCode string
	Summary    string
	Pos        token.Position
}

func getFuncSource(node *ast.FuncDecl, fset *token.FileSet) (string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, node)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Builds a call graph and optionally generates summaries
func BuildCodeFlowDiagram(node *ast.File, fset *token.FileSet, enableSummaries bool) (map[string]FunctionInfo, error) {
	funcMap := make(map[string]FunctionInfo)

	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			return true
		}
		funcName := fn.Name.Name

		src, err := getFuncSource(fn, fset)
		if err != nil {
			return true // skip errors
		}

		calls := []string{}
		ast.Inspect(fn.Body, func(bn ast.Node) bool {
			call, ok := bn.(*ast.CallExpr)
			if !ok {
				return true
			}
			switch fun := call.Fun.(type) {
			case *ast.Ident:
				calls = append(calls, fun.Name)
			case *ast.SelectorExpr:
				calls = append(calls, fun.Sel.Name)
			}
			return true
		})

		summary := ""
		if enableSummaries {
			s, err := api.GetFunctionSummary(src)
			if err == nil {
				summary = s
			}
		}

		funcMap[funcName] = FunctionInfo{
			Name:       funcName,
			Calls:      calls,
			SourceCode: src,
			Summary:    summary,
			Pos:        fset.Position(fn.Pos()),
		}

		return true
	})

	return funcMap, nil
}
