package providers

import (
	"fmt"
	"net/http"

	"gitlab.com/InfoBlogFriends/server/types"
	"gitlab.com/InfoBlogFriends/server/utils"

	"golang.org/x/oauth2/google"

	"gitlab.com/InfoBlogFriends/server/decoder"
	"gitlab.com/InfoBlogFriends/server/request"

	"golang.org/x/oauth2"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

type Google struct {
	*decoder.Decoder
	config *oauth2.Config
}

func NewGoogleAuth(config *oauth2.Config) *Google {
	config.Endpoint = google.Endpoint
	config.Scopes = []string{
		"https://www.googleapis.com/auth/userinfo.profile",
	}

	return &Google{config: config, Decoder: decoder.NewDecoder()}
}

func (f *Google) RedirectUrl() string {
	uuid, err := utils.ProjectUUIDGen("G")
	if err != nil {
		return ""
	}
	url := f.config.AuthCodeURL(uuid)

	return url
}

func (f *Google) Callback(r *http.Request) (infoblog.User, error) {
	code := r.FormValue("code")

	token, err := f.config.Exchange(r.Context(), code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		return infoblog.User{}, err
	}

	client := f.config.Client(r.Context(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		fmt.Printf("Get: %s\n", err)
		return infoblog.User{}, err
	}
	defer resp.Body.Close()

	var req request.GoogleCallbackResponse
	err = f.Decode(resp.Body, &req)
	if err != nil {
		return infoblog.User{}, err
	}

	return infoblog.User{
		GoogleID: types.NewNullString(req.GoogleID),
		Name:     types.NewNullString(req.Name),
	}, nil
}
