package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
	RedisURL    string
}

func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	fmt.Println("Using config file:", viper.ConfigFileUsed())

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		DatabaseURL: viper.GetString("DATABASE_URL"),
		ServerPort:  viper.GetString("SERVER_PORT"),
		RedisURL:    viper.GetString("REDIS_URL"),
	}, nil
}
