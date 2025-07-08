package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/wymersam/goflow/api"
)

func HandleSummaries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.FuncSummaries)
}
