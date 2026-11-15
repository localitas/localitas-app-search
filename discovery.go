package search

import (
	"encoding/json"
	"net/http"
)

type HealthInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Icon    string `json:"icon"`
	Status  string `json:"status"`
}

var DefaultHealth = HealthInfo{Name: "search", Version: "0.1.0", Icon: "search", Status: "healthy"}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DefaultHealth)
}
