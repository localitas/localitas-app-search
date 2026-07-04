package search

import (
	"net/http"

	"github.com/localitas/localitas-go/httputil"
)

type handler struct {
	searcher *Searcher
}

func (h *handler) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		writeErr(w, r, http.StatusBadRequest, "q is required")
		return
	}
	results := h.searcher.SearchAll(r.Context(), q)
	writeJSON(w, r, http.StatusOK, results)
}

func writeJSON(w http.ResponseWriter, r *http.Request, status int, v interface{}) {
	httputil.WriteResponse(w, r, status, v)
}

func writeErr(w http.ResponseWriter, r *http.Request, status int, format string, args ...interface{}) {
	httputil.WriteError(w, r, status, format, args...)
}
