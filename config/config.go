package config

import (
	"github.com/joho/godotenv"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	REST struct {
		Host string `default:"127.0.0.1" envconfig:"HOST"`
		Port string `default:"3000" envconfig:"PORT"`
	}
}

func Read() (Config, error) {
	var config Config

	err := godotenv.Load()
	if err != nil {
		return config, err
	}

	if err := envconfig.Process("WALLET_API", &config); err != nil {
		return config, err
	}

	return config, nil
}
