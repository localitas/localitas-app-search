package search

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleSwagger(t *testing.T) {
	w := httptest.NewRecorder()
	HandleSwagger(w, httptest.NewRequest("GET", "/swagger.json", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var spec APIDoc
	if err := json.Unmarshal(w.Body.Bytes(), &spec); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if spec.AppName != "Web Search" {
		t.Errorf("expected app_name Web Search, got %q", spec.AppName)
	}
	if len(spec.Endpoints) == 0 {
		t.Error("expected at least one endpoint")
	}
	hasSearch := false
	for _, ep := range spec.Endpoints {
		if strings.Contains(ep.Path, "/api/search") {
			hasSearch = true
		}
	}
	if !hasSearch {
		t.Error("expected /api/search endpoint")
	}
}

func TestHandleSearch_MissingQuery(t *testing.T) {
	h := &handler{searcher: NewSearcher(nil)}
	w := httptest.NewRecorder()
	h.handleSearch(w, httptest.NewRequest("GET", "/api/search", nil))
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestParseDuckDuckGoHTML(t *testing.T) {
	html := `
	<a class="result__a" href="//duckduckgo.com/l/?uddg=https%3A%2F%2Fgo.dev%2F&rut=abc">The Go Programming Language</a>
	<a class="result__snippet">Go is an open source programming language</a>
	<a class="result__a" href="//duckduckgo.com/l/?uddg=https%3A%2F%2Fen.wikipedia.org%2Fwiki%2FGo&rut=xyz">Go Wikipedia</a>
	<a class="result__snippet">Go is a programming language</a>
	`
	results := parseDuckDuckGoHTML(html)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Title != "The Go Programming Language" {
		t.Errorf("title = %s", results[0].Title)
	}
	if results[0].URL != "https://go.dev/" {
		t.Errorf("url = %s", results[0].URL)
	}
}

func TestParseBraveHTML(t *testing.T) {
	html := `
	<div class="snippet  svelte-jmfu5f" data-pos="1" data-type="web" data-keynav="true">
		<a href="https://go.dev/" target="_self" class="svelte-14r20fy l1">
			<div class="title search-snippet-title line-clamp-1 svelte-14r20fy" title="The Go Programming Language">The Go Programming Language</div>
		</a>
		<div class="content desktop-default-regular t-primary">Go is an open source programming language</div>
	</div>
	`
	results := parseBraveHTML(html)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].URL != "https://go.dev/" {
		t.Errorf("url = %s", results[0].URL)
	}
}

func TestParseBingHTML(t *testing.T) {
	html := `<li class="b_algo"><a href="https://go.dev/">The Go <strong>Programming</strong> Language</a><p class="desc">Go is an open source programming language</p></li>`
	results := parseBingHTML(html)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Title != "The Go Programming Language" {
		t.Errorf("title = %s", results[0].Title)
	}
}

func TestDecodeURL(t *testing.T) {
	encoded := "//duckduckgo.com/l/?uddg=https%3A%2F%2Fgo.dev%2F&rut=abc"
	decoded := decodeURL(encoded)
	if decoded != "https://go.dev/" {
		t.Errorf("expected https://go.dev/, got %s", decoded)
	}
}

func TestCleanHTML(t *testing.T) {
	html := "Hello <strong>World</strong> &amp; Friends"
	cleaned := cleanHTML(html)
	if cleaned != "Hello World & Friends" {
		t.Errorf("expected 'Hello World & Friends', got %s", cleaned)
	}
}

func TestSearchResults_EmptySlices(t *testing.T) {
	r := &SearchResults{
		Query:      "test",
		DuckDuckGo: make([]SearchResult, 0),
		Brave:      make([]SearchResult, 0),
		Bing:       make([]SearchResult, 0),
		Errors:     make([]string, 0),
	}
	b, _ := json.Marshal(r)
	s := string(b)
	if !contains(s, `"duckduckgo":[]`) {
		t.Error("expected empty array, not null")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
