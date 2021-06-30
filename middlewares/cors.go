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

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization")

		next.ServeHTTP(w, r)
	})
}
