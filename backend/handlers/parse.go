package handlers

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"net/http"
	"sort"
)

type ParseRequest struct {
	Path      string `json:"path"`
	EntryFunc string `json:"entryFunc"`
}

// Update response type to use Links instead of Edges
type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type ParseResponse struct {
	Nodes []string `json:"nodes"`
	Links []Link   `json:"links"`
}

func HandleParse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (max 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file in request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read uploaded file into memory
	srcBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read uploaded file", http.StatusInternalServerError)
		return
	}

	fileSet := token.NewFileSet()
	// Parse the file from bytes (instead of path)
	node, err := parser.ParseFile(fileSet, "", srcBytes, parser.AllErrors)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse Go source: %v", err), http.StatusBadRequest)
		return
	}

	// Build function map (assuming BuildCodeFlowDiagram supports this)
	funcMap, err := BuildCodeFlowDiagram(node, fileSet, true)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to build diagram: %v", err), http.StatusInternalServerError)
		return
	}

	// Collect nodes and links from funcMap as before
	nodeSet := make(map[string]struct{})
	var links []Link
	for funcName, info := range funcMap {
		nodeSet[funcName] = struct{}{}
		for _, call := range info.Calls {
			nodeSet[call] = struct{}{}
			links = append(links, Link{Source: funcName, Target: call})
		}
	}

	// Sort node list
	var nodes []string
	for n := range nodeSet {
		nodes = append(nodes, n)
	}
	sort.Strings(nodes)

	resp := ParseResponse{
		Nodes: nodes,
		Links: links,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
	}
}
