package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.com/InfoBlogFriends/server/types"

	"gitlab.com/InfoBlogFriends/server/auth"

	"gitlab.com/InfoBlogFriends/server/providers"

	"github.com/go-chi/chi/v5"

	"gitlab.com/InfoBlogFriends/server/cache"

	"gitlab.com/InfoBlogFriends/server/respond"
)

type AuthSocials struct {
	respond.Responder
	cache    cache.Cache
	facebook providers.Socials
	google   providers.Socials
}

func NewAuthSocials(responder respond.Responder, cache cache.Cache, facebook providers.Socials, google providers.Socials) *AuthSocials {
	return &AuthSocials{
		Responder: responder,
		cache:     cache,
		facebook:  facebook,
		google:    google,
	}
}

func (a *AuthSocials) Callback(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		providerName := chi.URLParam(r, "provider")
		var provider providers.Socials
		switch providerName {
		case "facebook":
			provider = a.facebook
		case "google":
			provider = a.google
		default:
			a.ErrorForbidden(w, fmt.Errorf("provider %s not exist", providerName))
			return
		}
		u, err := provider.Callback(r)
		if err != nil {
			a.ErrorForbidden(w, err)
			return
		}
		state := r.FormValue("state")
		ctx := context.WithValue(r.Context(), types.User{}, &u)
		ctx = context.WithValue(ctx, auth.State{}, state)
		ctx = context.WithValue(ctx, auth.Provider{}, providerName)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthSocials) Redirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		providerName := chi.URLParam(r, "provider")
		var provider providers.Socials
		switch providerName {
		case "facebook":
			provider = a.facebook
		case "google":
			provider = a.google
		default:
			a.ErrorForbidden(w, fmt.Errorf("provider %s not exist", providerName))
			return
		}
		u := provider.RedirectUrl()
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	})
}
