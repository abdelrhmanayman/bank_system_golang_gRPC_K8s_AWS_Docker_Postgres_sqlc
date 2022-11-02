package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port         string `mapstructure:"PORT"`
	DbDriver     string `mapstructure:"DB_DRIVER"`
	DbSourceName string `mapstructure:"DB_SOURCE_NAME"`
}

func LoadConfig() (config Config, err error) {
	viper.SetConfigFile(".env")
	err = viper.ReadInConfig()
	err = viper.Unmarshal(&config)
	return
}
