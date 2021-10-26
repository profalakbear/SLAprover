package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App AppConfig
}

type AppConfig struct {
	Host string
	Port string
	URN  string
}

func NewConfig() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
	config := &Config{
		App: AppConfig{
			Port: viper.Get("App.Port").(string),
			Host: viper.Get("App.Host").(string),
			URN:  viper.Get("App.Urn").(string),
		},
	}
	return config
}
