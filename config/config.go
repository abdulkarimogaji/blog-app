package config

import (
	"flag"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	PORT           string `mapstructure:"PORT"`
	DB_URI         string `mapstructure:"DB_URI"`
	JWT_SECRET     string `mapstructure:"JWT_SECRET"`
	GMAIL_NAME     string `mapstructure:"GMAIL_NAME"`
	GMAIL_ADDRESS  string `mapstructure:"GMAIL_ADDRESS"`
	GMAIL_PASSWORD string `mapstructure:"GMAIL_PASSWORD"`
	REDIS_ADDRESS  string `mapstructure:"REDIS_ADDRESS"`
}

var AppConfig Config

func LoadConfig() error {
	// init env variables
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	// read from .env
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// overwrite with command line flags

	// register flags here
	flag.String("PORT", "8080", "Port value")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	// bind to viper
	err = viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		return err
	}
	return nil
}
