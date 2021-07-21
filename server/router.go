package server

import (
	"bytes"
	"net/http"
	"time"

	"gitlab.com/InfoBlogFriends/server/components"
	"gitlab.com/InfoBlogFriends/server/services"

	"gitlab.com/InfoBlogFriends/server/email"
	"gitlab.com/InfoBlogFriends/server/request"

	"github.com/go-chi/cors"

	"gitlab.com/InfoBlogFriends/server/controllers"

	"gitlab.com/InfoBlogFriends/server/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(services services.Services, cmps components.Componenter) (*chi.Mux, error) {
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

	authController := controllers.NewAuth(cmps.Responder(), services.AuthService, cmps.Logger())

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		msg := email.NewMessage()
		var b bytes.Buffer

		b.Write([]byte("test"))

		msg.SetBody(b)
		msg.SetReceiver("globallinkliberty@gmail.com")
		msg.SetSubject("test")
		msg.OpenFile(".gitignore")
		msg.OpenFile(".env")
		err := cmps.Email().Send(msg)

		cmps.Responder().SendJSON(w, request.Response{
			Success: err == nil,
			Msg:     "test",
			Data:    err,
		})
	})

	r.Get("/swagger", swaggerUI)
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP(w, r)
	})
	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))).ServeHTTP(w, r)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/email/registration", authController.EmailActivation())
		r.Post("/email/verification", authController.EmailVerification())
		r.Post("/email/login", authController.EmailLogin())
		r.Post("/checkemail", authController.CheckCode())

		r.Post("/token/refresh", authController.RefreshToken())

		r.Post("/code", authController.SendCode())
		r.Post("/checkcode", authController.CheckCode())
	})

	token := middlewares.NewCheckToken(cmps.Responder(), cmps.JWTKeys())
	profileController := controllers.NewProfileController(cmps.Responder(), services.User, cmps.Logger())
	r.Route("/profile", func(r chi.Router) {
		r.Use(token.Check)
		r.Post("/update", profileController.Update())
		r.Patch("/update", profileController.Update())
		r.Get("/get", profileController.GetProfile())
		r.Route("/set", func(r chi.Router) {
			r.Post("/password", profileController.SetPassword())
		})
	})

	posts := controllers.NewPostsController(cmps.Responder(), services.User, services.File, services.Post, cmps.Logger())
	r.Route("/posts", func(r chi.Router) {
		r.Use(token.Check)
		r.Post("/like", posts.Like())
		r.Post("/create", posts.Create())
		r.Post("/update/{UUID}", posts.Create())
		r.Post("/delete/{UUID}", posts.Create())
		r.Get("/get/uuid/{UUID}", posts.Create())

		r.Route("/feed", func(r chi.Router) {
			r.Post("/my", posts.FeedMy())
			r.Post("/recent", posts.FeedRecent())

			r.Get("/subscribed", posts.FeedRecent())
			r.Post("/user", posts.FeedUser())
		})
	})

	users := controllers.NewUsersController(cmps.Responder(), services.User, cmps.Logger())
	r.Route("/user", func(r chi.Router) {
		r.Use(token.Check)
		r.Post("/subscribe", users.Subscribe())
		r.Post("/unsubscribe", users.Unsubscribe())
		r.Get("/list", users.List())
	})

	r.Route("/system", func(r chi.Router) {
		r.Use(middleware.Timeout(200 * time.Millisecond))
		r.Use(token.Check)
		r.Get("/config", func(w http.ResponseWriter, r *http.Request) {
			cmps.Responder().SendJSON(w, cmps.Config())
		})
	})

	return r, nil
}
