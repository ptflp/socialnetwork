package server

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func NewRouter(services *Services, components *HandlerComponents) (*chi.Mux, error) {
	r := chi.NewRouter()

	authHandler, err := NewAuthHandler(components.Responder, services.AuthService, components.Logger)

	r.Get("/swagger", swaggerUI)
	FileServer(r)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/code", authHandler.SendCode())
	})

	return r, err
}

func FileServer(r *chi.Mux) {
	root := "./static"
	fs := http.FileServer(http.Dir(root))

	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
			return
		}
		fs.ServeHTTP(w, r)
	})
}
