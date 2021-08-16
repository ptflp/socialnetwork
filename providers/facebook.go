package providers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"gitlab.com/InfoBlogFriends/server/types"

	"gitlab.com/InfoBlogFriends/server/decoder"
	"gitlab.com/InfoBlogFriends/server/request"

	"golang.org/x/oauth2/facebook"

	"gitlab.com/InfoBlogFriends/server/utils"

	"golang.org/x/oauth2"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

type Facebook struct {
	*decoder.Decoder
	config *oauth2.Config
}

func NewFacebookAuth(config *oauth2.Config) *Facebook {
	config.Endpoint = facebook.Endpoint
	config.Scopes = []string{"public_profile"}

	return &Facebook{config: config, Decoder: decoder.NewDecoder()}
}

func (f *Facebook) RedirectUrl() string {
	uuid, err := utils.ProjectUUIDGen("F")
	if err != nil {
		return ""
	}
	url := f.config.AuthCodeURL(uuid)

	return url
}

func (f *Facebook) Callback(r *http.Request) (infoblog.User, error) {
	code := r.FormValue("code")

	token, err := f.config.Exchange(r.Context(), code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		return infoblog.User{}, err
	}

	resp, err := http.Get("https://graph.facebook.com/me?access_token=" +
		url.QueryEscape(token.AccessToken))
	if err != nil {
		fmt.Printf("Get: %s\n", err)
		return infoblog.User{}, err
	}
	defer resp.Body.Close()

	var req request.FacebookCallbackRequest
	err = f.Decode(resp.Body, &req)
	if err != nil {
		return infoblog.User{}, err
	}
	facebookID, err := strconv.Atoi(req.FacebookID)
	if err != nil {
		return infoblog.User{}, err
	}

	return infoblog.User{
		FacebookID: types.NewNullInt64(int64(facebookID)),
		Name:       types.NewNullString(req.Name),
	}, nil
}
