package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App    App
	DB     DB
	Server Server
	Redis  Redis
}

func NewConfig() (*Config, error) {

	var config *Config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file, %s\n", err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("config unmarshall error")
	}

	return config, nil
}
