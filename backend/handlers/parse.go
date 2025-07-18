package handlers

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type ParseResponse struct {
	Nodes     []string          `json:"nodes"`
	Links     []Link            `json:"links"`
	Summaries map[string]string `json:"summaries"`
}

type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

func HandleRepoParse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(100 << 20) // 100MB max
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("repo")
	if err != nil {
		http.Error(w, "Missing repo file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	enableSummaries := r.FormValue("enableSummary") == "true"

	// Read ZIP into memory
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	if err != nil {
		http.Error(w, "Failed to read zip file", http.StatusInternalServerError)
		return
	}

	tmpDir, err := os.MkdirTemp("", "gorepo")
	if err != nil {
		http.Error(w, "Failed to create temp dir", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tmpDir)

	err = unzip(buf.Bytes(), tmpDir)
	if err != nil {
		http.Error(w, "Failed to unzip repo", http.StatusInternalServerError)
		return
	}

	// Parse all .go files
	fset := token.NewFileSet()
	var files []*ast.File

	err = filepath.WalkDir(tmpDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		node, err := parser.ParseFile(fset, path, src, parser.AllErrors)
		if err == nil {
			files = append(files, node)
		}
		return nil
	})
	if err != nil {
		http.Error(w, "Failed to walk files", http.StatusInternalServerError)
		return
	}

	funcMap, err := BuildCodeFlowDiagram(files, fset, enableSummaries)
	if err != nil {
		http.Error(w, "Failed to analyze code", http.StatusInternalServerError)
		return
	}

	// Build response
	summaries := make(map[string]string)
	nodeSet := make(map[string]struct{})
	var links []Link

	for funcName, info := range funcMap {
		summaries[funcName] = info.Summary
		nodeSet[funcName] = struct{}{}
		for _, call := range info.Calls {
			nodeSet[call] = struct{}{}
			links = append(links, Link{Source: funcName, Target: call})
		}
	}

	var nodes []string
	for n := range nodeSet {
		nodes = append(nodes, n)
	}
	sort.Strings(nodes)

	resp := ParseResponse{
		Nodes:     nodes,
		Links:     links,
		Summaries: summaries,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func unzip(data []byte, dest string) error {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}
	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}
		outFile, err := os.Create(fpath)
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
