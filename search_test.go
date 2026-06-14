package search

import (
	"testing"
)

func TestCleanHTML_Tags(t *testing.T) {
	input := "<b>Hello</b> <i>World</i>"
	got := cleanHTML(input)
	if got != "Hello World" {
		t.Errorf("expected 'Hello World', got %q", got)
	}
}

func TestCleanHTML_Entities(t *testing.T) {
	tests := map[string]string{
		"A &amp; B":          "A & B",
		"1 &lt; 2 &gt; 0":    "1 < 2 > 0",
		"&quot;quoted&quot;": "\"quoted\"",
		"it&#39;s":           "it's",
		"hello&nbsp;world":   "hello world",
	}
	for input, expected := range tests {
		got := cleanHTML(input)
		if got != expected {
			t.Errorf("cleanHTML(%q) = %q, want %q", input, got, expected)
		}
	}
}

func TestCleanHTML_Mixed(t *testing.T) {
	input := "  <p>Hello &amp; <b>World</b></p>  "
	got := cleanHTML(input)
	if got != "Hello & World" {
		t.Errorf("expected 'Hello & World', got %q", got)
	}
}

func TestDecodeURL_Regular(t *testing.T) {
	input := "https%3A%2F%2Fexample.com%2Fpath"
	got := decodeURL(input)
	if got != "https://example.com/path" {
		t.Errorf("expected decoded URL, got %q", got)
	}
}

func TestDecodeURL_DuckDuckGo_UddG(t *testing.T) {
	input := "/l/?kh=-1&uddg=https%3A%2F%2Fexample.com%2Fpage&rut=abc"
	got := decodeURL(input)
	if got != "https://example.com/page" {
		t.Errorf("expected 'https://example.com/page', got %q", got)
	}
}

func TestDecodeURL_AlreadyDecoded(t *testing.T) {
	input := "https://example.com/path"
	got := decodeURL(input)
	if got != "https://example.com/path" {
		t.Errorf("expected unchanged URL, got %q", got)
	}
}

func TestDecodeURL_Invalid(t *testing.T) {
	input := "%ZZinvalid"
	got := decodeURL(input)
	if got != input {
		t.Errorf("expected original string on error, got %q", got)
	}
}

func TestParseDuckDuckGoHTML_Empty(t *testing.T) {
	results := parseDuckDuckGoHTML("")
	if len(results) != 0 {
		t.Errorf("expected 0 results from empty HTML, got %d", len(results))
	}
}

func TestParseBraveHTML_Empty(t *testing.T) {
	results := parseBraveHTML("")
	if len(results) != 0 {
		t.Errorf("expected 0 results from empty HTML, got %d", len(results))
	}
}

func TestParseBingHTML_Empty(t *testing.T) {
	results := parseBingHTML("")
	if len(results) != 0 {
		t.Errorf("expected 0 results from empty HTML, got %d", len(results))
	}
}

func TestNewSearcher_Defaults(t *testing.T) {
	s := NewSearcher(nil)
	if s.client == nil {
		t.Error("expected non-nil client")
	}
}
