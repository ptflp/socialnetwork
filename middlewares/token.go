package middlewares

import (
	"context"
	"errors"
	"net/http"

	"gitlab.com/InfoBlogFriends/server/types"

	infoblog "gitlab.com/InfoBlogFriends/server"

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

func (t *Token) CheckStrict(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := t.jwt.ExtractAccessToken(r)
		if err != nil && (err.Error() == "token expired" || err.Error() == "Token is expired") {
			t.ErrorUnauthorized(w, errors.New("token expired"))
			return
		}
		if err != nil {
			t.ErrorForbidden(w, err)
			return
		}
		ctx := context.WithValue(r.Context(), types.User{}, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (t *Token) Check(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := t.jwt.ExtractAccessToken(r)
		if err != nil {
			u = &infoblog.User{}
		}
		ctx := context.WithValue(r.Context(), infoblog.User{}, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
