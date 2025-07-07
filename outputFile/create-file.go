package outputFile

import (
	"fmt"
	"os"
)

var CodeFlowGraph = make(map[string][]string)

func PrintMermaidToFile(fn string, visited map[string]bool, file *os.File, entryFunc string) {
	if visited[fn] {
		return
	}
	visited[fn] = true

	nodeClass := "normalFunc"
	callees := CodeFlowGraph[fn]

	if fn == entryFunc {
		nodeClass = "entryFunc"
	} else if len(callees) == 0 {
		nodeClass = "leafFunc"
	}

	fmt.Fprintf(file, "    %s[%q]:::%s\n", fn, fn, nodeClass)

	for _, callee := range callees {
		fmt.Fprintf(file, "    %s --> %s\n", fn, callee)
		PrintMermaidToFile(callee, visited, file, entryFunc)
	}
}
