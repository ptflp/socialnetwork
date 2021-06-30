package config

type DB struct {
	Net      string
	Driver   string
	DBName   string
	Username string `json:"-"`
	Password string `json:"-"`
	Host     string
	Port     string
	Timeout  int
}
