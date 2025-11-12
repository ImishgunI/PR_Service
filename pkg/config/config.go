package config

import (
	"errors"

	"github.com/spf13/viper"
)

func InitConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return errors.New("Reading .env file failed\n")
	}
	return nil
}

func GetString(key string) string {
	return viper.GetString(key)
}
