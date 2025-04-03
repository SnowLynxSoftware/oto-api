package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HealthController struct {
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (s *HealthController) MapController() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	return r
}
