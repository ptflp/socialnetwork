package server

import (
	"net/http"

	"github.com/go-chi/cors"

	"gitlab.com/InfoBlogFriends/server/config"

	"gitlab.com/InfoBlogFriends/server/handlers"

	"gitlab.com/InfoBlogFriends/server/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(services *Services, components *HandlerComponents, cfg *config.Config) (*chi.Mux, error) {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	authHandler := handlers.NewAuthHandler(components.Responder, services.AuthService, components.Logger)

	r.Get("/swagger", swaggerUI)
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP(w, r)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/code", authHandler.SendCode())
		r.Post("/checkcode", authHandler.CheckCode())
	})

	token := middlewares.NewCheckToken(components.Responder, components.JWTKeys)
	profileHandler := handlers.NewProfileHandler(components.Responder, services.User, components.Logger)
	r.Route("/profile", func(r chi.Router) {
		r.Use(token.Check)
		r.Post("/update", profileHandler.Update())
		r.Get("/get", profileHandler.GetProfile())
		r.Route("/set", func(r chi.Router) {
			r.Post("/password", profileHandler.SetPassword())
		})
	})

	r.Route("/system", func(r chi.Router) {
		r.Use(token.Check)
		r.Get("/config", func(w http.ResponseWriter, r *http.Request) {
			components.Responder.SendJSON(w, cfg)
		})
	})

	return r, nil
}
