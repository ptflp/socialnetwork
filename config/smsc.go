package config

type SMSC struct {
	Pwd   string `json:"-"`
	Login string `json:"-"`
	Cost  string
	Fmt   string
	Dev   bool
}
