package search

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

// An engine failing for a moment is asked again, with the browser headers
// it answers to, and its results come back with no error.
func TestSearchAll_AMomentaryEngineFailureIsAskedAgain(t *testing.T) {
	var mu sync.Mutex
	calls := map[string]int{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		calls[r.URL.Path]++
		n := calls[r.URL.Path]
		mu.Unlock()
		if n == 1 {
			http.Error(w, "down", http.StatusServiceUnavailable)
			return
		}
		if !strings.Contains(r.Header.Get("Accept"), "text/html") || r.Header.Get("Accept-Language") == "" || !strings.HasPrefix(r.UserAgent(), "Mozilla/") {
			http.Error(w, "not a browser", http.StatusBadRequest)
			return
		}
		if r.URL.Path == "/ddg" && r.URL.Query().Get("q") == "go" {
			w.Write([]byte(`<a class="result__a" href="https://go.dev/">The Go Programming Language</a><a class="result__snippet">Go is open source</a>`))
		}
	}))
	t.Cleanup(srv.Close)
	oldDDG, oldBrave, oldBing := duckDuckGoURL, braveURL, bingURL
	duckDuckGoURL, braveURL, bingURL = srv.URL+"/ddg", srv.URL+"/brave", srv.URL+"/bing"
	t.Cleanup(func() { duckDuckGoURL, braveURL, bingURL = oldDDG, oldBrave, oldBing })

	res := NewSearcher().SearchAll(context.Background(), "go")
	if len(res.Errors) != 0 || len(res.DuckDuckGo) != 1 || res.DuckDuckGo[0].URL != "https://go.dev/" {
		t.Errorf("results = %+v", res)
	}
	for _, p := range []string{"/ddg", "/brave", "/bing"} {
		if calls[p] != 2 {
			t.Errorf("%s asked %d times, want 2", p, calls[p])
		}
	}
}
