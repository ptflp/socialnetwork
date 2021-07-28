package providers

import (
	"net/http"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

type Socials interface {
	Callback(r *http.Request) (infoblog.User, error)
	RedirectUrl() string
}
