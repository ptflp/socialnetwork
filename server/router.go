package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(services *Services, components *HandlerComponents) (*chi.Mux, error) {
	r := chi.NewRouter()

	authHandler, err := NewAuthHandler(components.Responder, services.AuthService, components.Logger)

	r.Get("/swagger", swaggerUI)
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP(w, r)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/code", authHandler.SendCode())
		r.Post("/checkcode", authHandler.CheckCode())
	})

	return r, err
}
