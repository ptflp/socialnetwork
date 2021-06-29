package middlewares

import (
	"context"
	"net/http"

	"gitlab.com/InfoBlogFriends/server/session"

	"gitlab.com/InfoBlogFriends/server/respond"
)

type Token struct {
	respond.Responder
	jwt *session.JWTKeys
}

func NewCheckToken(responder respond.Responder, jwt *session.JWTKeys) *Token {
	return &Token{
		Responder: responder,
		jwt:       jwt,
	}
}

func (t *Token) Check(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := t.jwt.ExtractToken(r)
		if err != nil {
			t.ErrorForbidden(w, err)
			return
		}
		ctx := context.WithValue(r.Context(), "user", u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
