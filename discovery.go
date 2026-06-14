package search

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/grandcat/zeroconf"
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

func BroadcastMDNS(port int, name string) (func(), error) {
	server, err := zeroconf.Register(name, "_localitas-app._tcp", "local.", port, []string{"name=" + name}, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("mDNS broadcasting %s on port %d", name, port)
	return func() { server.Shutdown() }, nil
}
