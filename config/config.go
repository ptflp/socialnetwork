package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	App    App
	DB     DB
	Server Server
	Redis  Redis
	SMSC   SMSC
	Email  Email
	Oauth2
}

const (
	ProductionKey = "production"
	DevKey        = "dev"
	StageKey      = "stage"
	Type          = "yaml"
	Path          = "./config"

	CheckEnvKey = "ENV"
)

func NewConfig() (*Config, error) {

	var config *Config

	viper.SetConfigName(ProductionKey)
	viper.SetConfigType(Type)
	viper.AddConfigPath(Path)

	v := viper.New()
	v.SetConfigName(os.Getenv(CheckEnvKey))
	v.SetConfigType(Type)
	v.AddConfigPath(Path)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file, config %s, %s", os.Getenv(CheckEnvKey), err)
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file, %s\n", err)
	}

	if os.Getenv(CheckEnvKey) != "production" {
		if err := viper.MergeConfigMap(v.AllSettings()); err != nil {
			return nil, fmt.Errorf("error merge dev config file, %s\n", err)
		}
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("config unmarshall error")
	}

	return config, nil
}
