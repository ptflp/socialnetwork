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
}

const (
	ProductionConfigName = "prod"
	DevConfigName        = "dev"
	Type                 = "yaml"
	Path                 = "./config"

	CheckEnvKey = "DEV"
)

func NewConfig() (*Config, error) {

	var config *Config

	viper.SetConfigName(ProductionConfigName)
	viper.SetConfigType(Type)
	viper.AddConfigPath(Path)

	v := viper.New()
	if os.Getenv(CheckEnvKey) == "true" {
		v.SetConfigName(DevConfigName)
		v.SetConfigType(Type)
		v.AddConfigPath(Path)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("error reading config file, %s\n", err)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file, %s\n", err)
	}

	if os.Getenv(CheckEnvKey) == "true" {
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
