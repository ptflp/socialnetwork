package middlewares

import (
	"net/http"

	"github.com/rs/cors"
)

type Cors struct {
}

func NewCors() *Cors {
	return &Cors{}
}

func (t *Cors) OpenAllCors(next http.Handler) http.Handler {

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Accept", "Accept-Language"},
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
		AllowCredentials: true,
		Debug:            false,
	})

	return cors.Handler(func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}())
}
