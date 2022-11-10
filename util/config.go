package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port          string        `mapstructure:"PORT"`
	DbDriver      string        `mapstructure:"DB_DRIVER"`
	DbSourceName  string        `mapstructure:"DB_SOURCE_NAME"`
	SymmetricKey  string        `mapstructure:"SYMMETRIC_KEY"`
	TokenDuration time.Duration `mapstructure:"TOKEN_DURATION"`
}

func LoadConfig() (config Config, err error) {
	viper.SetConfigFile(".env")
	err = viper.ReadInConfig()
	err = viper.Unmarshal(&config)
	return
}
