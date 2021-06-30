package middlewares

import (
	"net/http"
)

type Cors struct {
}

func NewCors() *Cors {
	return &Cors{}
}

func (t *Cors) OpenAllCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}
