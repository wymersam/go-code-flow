package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/wymersam/goflow/api"
	"github.com/wymersam/goflow/diagram"
	"github.com/wymersam/goflow/outputFile"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Loaded API Key:", os.Getenv("OPENAI_API_KEY"))
}

func main() {

	var generateSummaries bool
	flag.BoolVar(&generateSummaries, "summaries", false, "Generate function summaries (uses OpenAI)")
	flag.Parse()

	diagram.EnableSummaries = generateSummaries

	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Usage: go run main.go [flags] [path-to-code] [entry-func-name]")
		return
	}
	root := args[0]
	entryFunc := args[1]

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
		diagram.BuildCodeFlowDiagram(node, fileSet)
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
	outputFile.PrintMermaidToFile(entryFunc, visited, file, entryFunc)
	fmt.Fprintln(file, "```")
	fmt.Println("✅ Mermaid diagram written to codeflow.md")

	// Append Function Summaries
	fmt.Fprintln(file, "\n## 📘 Function Summaries\n")
	for fn, summary := range api.FuncSummaries {
		fmt.Fprintf(file, "<details>\n<summary><strong>%s</strong></summary>\n\n", fn)
		fmt.Fprintln(file, "```go")
		fmt.Fprintln(file, summary)
		fmt.Fprintln(file, "```\n</details>\n")
	}

	fmt.Println("✅ Function summaries written to codeflow.md")
}
