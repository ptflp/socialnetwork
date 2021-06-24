package config

type DB struct {
	Net      string
	Driver   string
	DBName   string
	Username string
	Password string
	Host     string
	Port     string
	Timeout  int
}
