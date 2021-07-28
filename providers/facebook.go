package providers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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
	config.RedirectURL = "https://cbf8129b5d2f.ngrok.io/auth/provider/facebook/callback"

	return &Facebook{config: config, Decoder: decoder.NewDecoder()}
}

func (f *Facebook) RedirectUrl() string {
	Url, err := url.Parse(f.config.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}
	uuid, err := utils.ProjectUUIDGen("F")
	if err != nil {
		return ""
	}
	parameters := url.Values{}
	parameters.Add("client_id", f.config.ClientID)
	parameters.Add("scope", strings.Join(f.config.Scopes, " "))
	parameters.Add("redirect_uri", f.config.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", uuid)
	Url.RawQuery = parameters.Encode()
	u := Url.String()

	return u
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
		FacebookID: infoblog.NewNullInt64(int64(facebookID)),
	}, nil
}
