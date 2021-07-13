package controllers

import (
	"errors"
	"net/http"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

func extractUser(r *http.Request) (infoblog.User, error) {
	ctx := r.Context()
	u, ok := ctx.Value("user").(*infoblog.User)
	if !ok {
		return infoblog.User{}, errors.New("type assertion to user err")
	}

	if u.ID == 0 {
		return infoblog.User{}, errors.New("user not exists")
	}

	return *u, nil
}
