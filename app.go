package search

import (
	"html/template"
	"log"
	"net/http"
)

type App struct {
	Searcher *Searcher
	BasePath string
}

func New(basePath string) *App {
	if basePath == "" {
		basePath = "/"
	}
	return &App{
		Searcher: NewSearcher(nil),
		BasePath: basePath,
	}
}

func (a *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(TemplatesFS, "templates/index.html")
	if err != nil {
		log.Printf("search index template error: %v", err)
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	tmpl.ExecuteTemplate(w, "index.html", map[string]string{
		"BasePath": a.BasePath,
	})
}

func (a *App) RegisterRoutes(mux *http.ServeMux) {
	h := &handler{searcher: a.Searcher}

	mux.HandleFunc("GET /{$}", a.handleIndex)
	mux.HandleFunc("GET /swagger.json", HandleSwagger)
	mux.HandleFunc("GET /help.md", handleHelpMarkdown)
	mux.HandleFunc("GET /api/search", h.handleSearch)
}
