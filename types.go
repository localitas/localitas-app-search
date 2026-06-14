package search

type SearchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

type SearchResults struct {
	Query      string         `json:"query"`
	DuckDuckGo []SearchResult `json:"duckduckgo"`
	Brave      []SearchResult `json:"brave"`
	Bing       []SearchResult `json:"bing"`
	Errors     []string       `json:"errors,omitempty"`
}
