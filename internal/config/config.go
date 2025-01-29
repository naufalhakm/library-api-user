package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost         string `mapstructure:"DB_HOST"`
	DBUserName     string `mapstructure:"DB_USERNAME"`
	DBUserPassword string `mapstructure:"DB_PASSWORD"`
	DBName         string `mapstructure:"DB_DATABASE"`
	DBPort         string `mapstructure:"DB_PORT"`
	ServerPort     string `mapstructure:"PORT"`
	GRPCPort       string `mapstructure:"GRCP_PORT"`
	BookGRCP       string `mapstructure:"BOOK_GRCP"`
}

var ENV *Config

func LoadConfig() {
	fang := viper.New()

	fang.AddConfigPath(".")
	fang.SetConfigName(".env")
	fang.SetConfigType("env")

	err := fang.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = fang.Unmarshal(&ENV)
	if err != nil {
		panic(err)
	}
}
