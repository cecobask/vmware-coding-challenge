package pagedata

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

// NewRouter creates all routes associated with accounts
func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", h.GetPageData)
	})
	return r
}
