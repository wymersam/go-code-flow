package diagram

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/wymersam/goflow/api"
	"github.com/wymersam/goflow/outputFile"
)

func BuildCodeFlowDiagram(node *ast.File, fset *token.FileSet) {
	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			return true
		}
		funcName := fn.Name.Name
		fmt.Println("Found function:", funcName)

		// Get the source code of this function
		src, err := outputFile.GetFuncSource(fn, fset)
		if err != nil {
			fmt.Println("Error getting source for", funcName, ":", err)
			return true
		}

		// Get the summary of the function
		summary, err := api.GetFunctionSummary(src)
		if err != nil {
			fmt.Println("Error getting summary for", funcName, ":", err)
			return true
		}

		api.FuncSummaries[funcName] = summary
		fmt.Printf("Summary for %s:\n%s\n\n", funcName, summary)

		// Build the code flow graph
		ast.Inspect(fn.Body, func(bn ast.Node) bool {
			call, ok := bn.(*ast.CallExpr)
			if !ok {
				return true
			}
			switch fun := call.Fun.(type) {
			case *ast.Ident:
				outputFile.CodeFlowGraph[funcName] = append(outputFile.CodeFlowGraph[funcName], fun.Name)
			case *ast.SelectorExpr:
				outputFile.CodeFlowGraph[funcName] = append(outputFile.CodeFlowGraph[funcName], fun.Sel.Name)
			}
			return true
		})
		return true
	})
}
