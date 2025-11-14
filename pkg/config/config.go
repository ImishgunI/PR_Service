package config

import (
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.ReadInConfig()
}

func GetString(key string) string {
	return viper.GetString(key)
}
