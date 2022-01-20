package handlers

import (
	"errors"
	"net/http"

	"gitlab.com/InfoBlogFriends/server/types"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

func extractUser(r *http.Request) (infoblog.User, error) {
	ctx := r.Context()
	u, ok := ctx.Value(types.User{}).(*infoblog.User)
	if !ok {
		return infoblog.User{}, errors.New("type assertion to user err")
	}

	return *u, nil
}
