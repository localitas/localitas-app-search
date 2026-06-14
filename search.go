package search

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

type Searcher struct {
	client *http.Client
}

func NewSearcher(client *http.Client) *Searcher {
	if client == nil {
		client = http.DefaultClient
	}
	return &Searcher{client: client}
}

func (s *Searcher) SearchAll(ctx context.Context, query string) *SearchResults {
	results := &SearchResults{
		Query:      query,
		DuckDuckGo: make([]SearchResult, 0),
		Brave:      make([]SearchResult, 0),
		Bing:       make([]SearchResult, 0),
		Errors:     make([]string, 0),
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(3)

	go func() {
		defer wg.Done()
		r, err := s.searchDuckDuckGo(ctx, query)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			results.Errors = append(results.Errors, fmt.Sprintf("DuckDuckGo: %v", err))
		} else {
			results.DuckDuckGo = r
		}
	}()

	go func() {
		defer wg.Done()
		r, err := s.searchBrave(ctx, query)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			results.Errors = append(results.Errors, fmt.Sprintf("Brave: %v", err))
		} else {
			results.Brave = r
		}
	}()

	go func() {
		defer wg.Done()
		r, err := s.searchBing(ctx, query)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			results.Errors = append(results.Errors, fmt.Sprintf("Bing: %v", err))
		} else {
			results.Bing = r
		}
	}()

	wg.Wait()
	return results
}

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

func (s *Searcher) fetch(ctx context.Context, searchURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (s *Searcher) searchDuckDuckGo(ctx context.Context, query string) ([]SearchResult, error) {
	body, err := s.fetch(ctx, fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", url.QueryEscape(query)))
	if err != nil {
		return nil, err
	}
	return parseDuckDuckGoHTML(body), nil
}

func (s *Searcher) searchBrave(ctx context.Context, query string) ([]SearchResult, error) {
	body, err := s.fetch(ctx, fmt.Sprintf("https://search.brave.com/search?q=%s", url.QueryEscape(query)))
	if err != nil {
		return nil, err
	}
	return parseBraveHTML(body), nil
}

func (s *Searcher) searchBing(ctx context.Context, query string) ([]SearchResult, error) {
	body, err := s.fetch(ctx, fmt.Sprintf("https://www.bing.com/search?q=%s", url.QueryEscape(query)))
	if err != nil {
		return nil, err
	}
	return parseBingHTML(body), nil
}

func parseDuckDuckGoHTML(html string) []SearchResult {
	results := make([]SearchResult, 0)
	resultPattern := regexp.MustCompile(`<a[^>]*class="result__a"[^>]*href="([^"]*)"[^>]*>([^<]*)</a>`)
	snippetPattern := regexp.MustCompile(`<a[^>]*class="result__snippet"[^>]*>([^<]*(?:<[^>]*>[^<]*</[^>]*>)*[^<]*)</a>`)

	titleMatches := resultPattern.FindAllStringSubmatch(html, -1)
	snippetMatches := snippetPattern.FindAllStringSubmatch(html, -1)

	for i, match := range titleMatches {
		if len(match) >= 3 {
			r := SearchResult{URL: decodeURL(match[1]), Title: cleanHTML(match[2])}
			if i < len(snippetMatches) && len(snippetMatches[i]) >= 2 {
				r.Snippet = cleanHTML(snippetMatches[i][1])
			}
			if r.URL != "" && r.Title != "" {
				results = append(results, r)
			}
		}
		if len(results) >= 10 {
			break
		}
	}
	return results
}

func parseBraveHTML(html string) []SearchResult {
	results := make([]SearchResult, 0)
	blockPattern := regexp.MustCompile(`<div[^>]*class="snippet[^"]*svelte[^"]*"[^>]*data-pos="\d+"[^>]*data-type="web"[^>]*>`)
	blockIndices := blockPattern.FindAllStringIndex(html, -1)

	for i, idx := range blockIndices {
		start := idx[0]
		end := len(html)
		if i+1 < len(blockIndices) {
			end = blockIndices[i+1][0]
		} else if start+5000 < len(html) {
			end = start + 5000
		}
		block := html[start:end]

		urlPattern := regexp.MustCompile(`<a[^>]*href="(https?://[^"]*)"[^>]*target="_self"[^>]*class="svelte[^"]*"[^>]*>`)
		titlePattern := regexp.MustCompile(`<div[^>]*class="title search-snippet-title[^"]*"[^>]*title="([^"]*)"[^>]*>`)
		snippetPattern := regexp.MustCompile(`<div[^>]*class="content desktop-default-regular[^"]*"[^>]*>([^<]*(?:<[^>]*>[^<]*</[^>]*>)*[^<]*)</div>`)

		urlMatch := urlPattern.FindStringSubmatch(block)
		titleMatch := titlePattern.FindStringSubmatch(block)
		snippetMatch := snippetPattern.FindStringSubmatch(block)

		if len(urlMatch) >= 2 && len(titleMatch) >= 2 {
			if strings.Contains(urlMatch[1], "brave.com") || strings.Contains(urlMatch[1], "search.brave") {
				continue
			}
			r := SearchResult{URL: urlMatch[1], Title: cleanHTML(titleMatch[1])}
			if len(snippetMatch) >= 2 {
				r.Snippet = cleanHTML(snippetMatch[1])
			}
			if r.Title != "" {
				results = append(results, r)
			}
		}
		if len(results) >= 10 {
			break
		}
	}
	return results
}

func parseBingHTML(html string) []SearchResult {
	results := make([]SearchResult, 0)
	resultPattern := regexp.MustCompile(`<li[^>]*class="b_algo"[^>]*>(.*?)</li>`)
	titlePattern := regexp.MustCompile(`<a[^>]*href="(https?://[^"]*)"[^>]*>([^<]*(?:<strong>[^<]*</strong>[^<]*)*)</a>`)
	snippetPattern := regexp.MustCompile(`<p[^>]*class="[^"]*"[^>]*>([^<]*(?:<[^>]*>[^<]*</[^>]*>)*[^<]*)</p>`)

	blocks := resultPattern.FindAllStringSubmatch(html, -1)
	for _, block := range blocks {
		titleMatch := titlePattern.FindStringSubmatch(block[1])
		snippetMatch := snippetPattern.FindStringSubmatch(block[1])
		if len(titleMatch) >= 3 {
			r := SearchResult{URL: titleMatch[1], Title: cleanHTML(titleMatch[2])}
			if len(snippetMatch) >= 2 {
				r.Snippet = cleanHTML(snippetMatch[1])
			}
			if r.Title != "" {
				results = append(results, r)
			}
		}
		if len(results) >= 10 {
			break
		}
	}
	return results
}

func decodeURL(encodedURL string) string {
	if strings.Contains(encodedURL, "uddg=") {
		parts := strings.Split(encodedURL, "uddg=")
		if len(parts) > 1 {
			decoded, err := url.QueryUnescape(parts[1])
			if err == nil {
				if idx := strings.Index(decoded, "&"); idx > 0 {
					decoded = decoded[:idx]
				}
				return decoded
			}
		}
	}
	decoded, err := url.QueryUnescape(encodedURL)
	if err != nil {
		return encodedURL
	}
	return decoded
}

func cleanHTML(s string) string {
	tagPattern := regexp.MustCompile(`<[^>]*>`)
	s = tagPattern.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	return strings.TrimSpace(s)
}
