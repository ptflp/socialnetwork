package controllers

import (
	"errors"
	"net/http"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

func extractUser(r *http.Request) (infoblog.User, error) {
	ctx := r.Context()
	u, ok := ctx.Value(infoblog.User{}).(*infoblog.User)
	if !ok {
		return infoblog.User{}, errors.New("type assertion to user err")
	}

	return *u, nil
}
