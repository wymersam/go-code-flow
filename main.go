package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

var codeFlowGraph = make(map[string][]string)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go [path-to-code] [entry-func-name]")
		return
	}
	root := os.Args[1]
	entryFunc := os.Args[2]

	// ✅ Check if path exists and is a directory
	info, err := os.Stat(root)
	if os.IsNotExist(err) {
		fmt.Printf("❌ Error: The path '%s' does not exist.\n", root)
		return
	}
	if !info.IsDir() {
		fmt.Printf("❌ Error: The path '%s' is not a directory.\n", root)
		return
	}

	fileSet := token.NewFileSet()

	err = filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if filepath.Ext(path) != ".go" {
			return nil
		}
		fmt.Println("Walking:", path)
		node, err := parser.ParseFile(fileSet, path, nil, parser.AllErrors)
		if err != nil {
			fmt.Println("Error parsing:", path, err)
			return nil
		}
		buildCodeFlowDiagram(node)
		return nil
	})
	if err != nil {
		fmt.Println("Error walking through code:", err)
		return
	}

	// Create Markdown file
	file, err := os.Create("codeflow.md")
	if err != nil {
		fmt.Println("Error creating Markdown file:", err)
		return
	}
	defer file.Close()

	// Write Mermaid content to file
	fmt.Fprintln(file, "# Function Call Graph\n")
	fmt.Fprintln(file, "```mermaid")

	// === ADD: change layout to LR and define a class for bigger font size ===
	fmt.Fprintln(file, "graph LR")
	fmt.Fprintln(file, "classDef bigFont fill:#fff,stroke:#333,stroke-width:1px,font-size:16px;")

	visited := make(map[string]bool)
	printMermaidToFile(entryFunc, visited, file)

	fmt.Fprintln(file, "```")

	fmt.Println("✅ Mermaid diagram written to codeflow.md")
}

func buildCodeFlowDiagram(node *ast.File) {
	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			return true
		}
		funcName := fn.Name.Name
		fmt.Println("Found function:", funcName)

		ast.Inspect(fn.Body, func(bn ast.Node) bool {
			call, ok := bn.(*ast.CallExpr)
			if !ok {
				return true
			}
			switch fun := call.Fun.(type) {
			case *ast.Ident:
				codeFlowGraph[funcName] = append(codeFlowGraph[funcName], fun.Name)
			case *ast.SelectorExpr:
				codeFlowGraph[funcName] = append(codeFlowGraph[funcName], fun.Sel.Name)
			}

			return true
		})
		return true
	})
}

func printMermaidToFile(fn string, visited map[string]bool, file *os.File) {
	if visited[fn] {
		return
	}
	visited[fn] = true

	callees := codeFlowGraph[fn]
	if len(callees) == 0 {
		// Leaf node: print node with bigFont class
		fmt.Fprintf(file, "    %s:::bigFont\n", fn)
		return
	}

	for _, callee := range callees {
		// Print edges without styling
		fmt.Fprintf(file, "    %s --> %s\n", fn, callee)
		// Print callee node with class
		fmt.Fprintf(file, "    %s:::bigFont\n", callee)
		printMermaidToFile(callee, visited, file)
	}
}
