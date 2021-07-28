package config

import "golang.org/x/oauth2"

type Oauth2 struct {
	Google   oauth2.Config
	Facebook oauth2.Config
}
