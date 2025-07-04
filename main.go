package main

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

var codeFlowGraph = make(map[string][]string)
var funcSummaries = make(map[string]string)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Loaded API Key:", os.Getenv("OPENAI_API_KEY"))
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go [path-to-code] [entry-func-name]")
		return
	}
	root := os.Args[1]
	entryFunc := os.Args[2]

	// Check if path exists and is a directory
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
		buildCodeFlowDiagram(node, fileSet)
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
	fmt.Fprintln(file, "# 🔁 Function Call Graph\n")
	fmt.Fprintln(file, "```mermaid")
	fmt.Fprintln(file, "graph LR")
	fmt.Fprintln(file, "classDef entryFunc fill:#f96,stroke:#333,stroke-width:2px,font-weight:bold,font-size:18px,color:#000;")
	fmt.Fprintln(file, "classDef leafFunc fill:#6f9,stroke:#333,stroke-width:1px,font-style:italic,font-size:14px,color:#000;")
	fmt.Fprintln(file, "classDef normalFunc fill:#fff,stroke:#333,stroke-width:1px,font-size:16px,color:#000;")

	visited := make(map[string]bool)
	printMermaidToFile(entryFunc, visited, file, entryFunc)
	fmt.Fprintln(file, "```")
	fmt.Println("✅ Mermaid diagram written to codeflow.md")

	// Append Function Summaries
	fmt.Fprintln(file, "\n## 📘 Function Summaries\n")
	for fn, summary := range funcSummaries {
		fmt.Fprintf(file, "<details>\n<summary><strong>%s</strong></summary>\n\n", fn)
		fmt.Fprintln(file, "```go")
		fmt.Fprintln(file, summary)
		fmt.Fprintln(file, "```\n</details>\n")
	}

	fmt.Println("✅ Function summaries written to codeflow.md")
}

func buildCodeFlowDiagram(node *ast.File, fset *token.FileSet) {
	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			return true
		}
		funcName := fn.Name.Name
		fmt.Println("Found function:", funcName)

		// Get the source code of this function
		src, err := getFuncSource(fn, fset)
		if err != nil {
			fmt.Println("Error getting source for", funcName, ":", err)
			return true
		}

		// Get the summary of the function
		summary, err := getFunctionSummary(src)
		if err != nil {
			fmt.Println("Error getting summary for", funcName, ":", err)
			return true
		}

		funcSummaries[funcName] = summary
		fmt.Printf("Summary for %s:\n%s\n\n", funcName, summary)

		// Build the code flow graph
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

func printMermaidToFile(fn string, visited map[string]bool, file *os.File, entryFunc string) {
	if visited[fn] {
		return
	}
	visited[fn] = true

	nodeClass := "normalFunc"
	callees := codeFlowGraph[fn]

	if fn == entryFunc {
		nodeClass = "entryFunc"
	} else if len(callees) == 0 {
		nodeClass = "leafFunc"
	}

	fmt.Fprintf(file, "    %s[%q]:::%s\n", fn, fn, nodeClass)

	for _, callee := range callees {
		fmt.Fprintf(file, "    %s --> %s\n", fn, callee)
		printMermaidToFile(callee, visited, file, entryFunc)
	}
}

func getFuncSource(node *ast.FuncDecl, fset *token.FileSet) (string, error) {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, node)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func getFunctionSummary(code string) (string, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "gpt-4o-mini",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "system",
					Content: "You are a helpful assistant. Summarise the following Go function in one or two sentences.",
				},
				{
					Role:    "user",
					Content: code,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}
