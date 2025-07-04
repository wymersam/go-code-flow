# Code Flow Generator (Go)

This tool parses a Go project directory and generates a Mermaid.js function call graph starting from a specified entry function. The output is saved as a Markdown file (`codeflow.md`) containing a Mermaid diagram that visualises the call relationships.

---

## Prerequisites

The following mermaid extension is needed to see the finished diagrams in VS Code:

```shell
code --install-extension bierner.markdown-mermaid
```

---

## Features

- Recursively parses all `.go` files in a given directory.
- Extracts function calls inside functions to build a code flow graph.
- Generates Mermaid syntax for a **left-to-right (LR)** graph layout.
- Highlights nodes with a larger font size for better readability.
- Allows you to specify the starting function to focus the graph on a particular use case or scenario.

---

## Usage

```bash
go run main.go [path-to-go-code] [entry-function-name]
```

**path-to-go-code**: Path to the root directory of your Go project.

**entry-function-name**: The function where the call graph traversal begins.

Example:

```bash
go run main.go ./myproject HandleRequest
```

This generates a file `codeflow.md` containing a Mermaid graph of calls starting from HandleRequest.

Output:

`codeflow.md` — Markdown file with a Mermaid diagram.

Example snippet inside codeflow.md:

```markdown
mermaid
Copy
Edit
graph LR
classDef bigFont fill:#fff,stroke:#333,stroke-width:1px,font-size:16px;

HandleRequest --> AuthenticateUser
AuthenticateUser:::bigFont
```

You can render this Markdown with the previosuly installed Mermaid extension or Mermaid live editors to visualise your function call graph.

### How It Works

The program walks through the Go source files in the given directory.

Parses each file using Go's `go/ast` and `go/parser` packages.

Builds a map of functions and their called functions.

Recursively writes Mermaid edges starting from the specified entry function.

Adds CSS classes to nodes for better font sizing.

Requirements
Go 1.16+

Compatible Mermaid renderer to visualize the diagram.

License
MIT License — Feel free to modify and use it in your projects.
