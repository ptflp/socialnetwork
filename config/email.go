package config

type Email struct {
	ServerAddress string
	Port          string
	Login         string `json:"-"`
	Password      string `json:"-"`
	From          string
}
