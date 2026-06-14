package search

import (
	"encoding/json"
	"net/http"
)

type APIEndpoint struct {
	Method      string     `json:"method"`
	Path        string     `json:"path"`
	Summary     string     `json:"summary"`
	QueryParams []APIParam `json:"query_params,omitempty"`
	Response    *APIBody   `json:"response,omitempty"`
}

type APIParam struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

type APIBody struct {
	ContentType string `json:"content_type"`
	Example     string `json:"example"`
}

type APIDoc struct {
	AppName     string        `json:"app_name"`
	Version     string        `json:"version"`
	Description string        `json:"description"`
	Keywords    []string      `json:"keywords,omitempty"`
	Endpoints   []APIEndpoint `json:"endpoints"`
}

var SearchAPIDoc = APIDoc{
	AppName:     "Web Search",
	Version:     "0.1.0",
	Description: "Parallel web search across DuckDuckGo, Brave, and Bing",
	Keywords:    []string{"search", "web", "google", "query", "find", "lookup", "browse", "internet", "results"},
	Endpoints: []APIEndpoint{
		{Method: "GET", Path: "/api/search", Summary: "Search the web across three engines", QueryParams: []APIParam{{Name: "q", Type: "string", Required: true, Description: "Search query"}}, Response: &APIBody{ContentType: "application/json", Example: `{"query":"golang","duckduckgo":[{"title":"Go","url":"https://go.dev/","snippet":"..."}],"brave":[],"bing":[]}`}},
	},
}

func HandleSwagger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SearchAPIDoc)
}
