package search

import (
	"net/http"

	client "github.com/localitas/localitas-go"
)

var SearchAPIDoc = client.APISpec{
	AppName:     "Web Search",
	Version:     "0.1.0",
	Description: "Parallel web search across DuckDuckGo, Brave, and Bing",
	Keywords:    []string{"search", "web", "google", "query", "find", "lookup", "browse", "internet", "results"},
	Endpoints: []client.APIEndpoint{
		{Method: "GET", Path: "/api/search", Summary: "Search the web across three engines", QueryParams: []client.APIParam{{Name: "q", Type: "string", Required: true, Description: "Search query"}}, Response: &client.APIBody{ContentType: "application/json", Example: `{"query":"golang","duckduckgo":[{"title":"Go","url":"https://go.dev/","snippet":"..."}],"brave":[],"bing":[]}`}},
	},
}

func HandleSwagger(w http.ResponseWriter, r *http.Request) {
	SearchAPIDoc.ServeHTTP(w, r)
}
