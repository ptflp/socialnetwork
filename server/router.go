package server

import (
	"errors"
	"net/http"

	infoblog "gitlab.com/InfoBlogFriends/server"

	"gitlab.com/InfoBlogFriends/server/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(services *Services, components *HandlerComponents) (*chi.Mux, error) {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	authHandler, err := NewAuthHandler(components.Responder, services.AuthService, components.Logger)

	r.Get("/swagger", swaggerUI)
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP(w, r)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/code", authHandler.SendCode())
		r.Post("/checkcode", authHandler.CheckCode())
	})

	token := middlewares.NewCheckToken(components.Responder, components.JWTKeys)
	r.Route("/test", func(r chi.Router) {
		r.Use(token.Check)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			u, ok := ctx.Value("user").(*infoblog.User)
			if !ok {
				components.Responder.ErrorInternal(w, errors.New("type assertion to user err"))
				return
			}
			components.Responder.SendJSON(w, u)
		})
	})

	return r, err
}