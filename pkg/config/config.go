package config

import (
	"PullRequestService/pkg/logger"

	"github.com/spf13/viper"
)

func InitConfig() {
	log := logger.New()
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Warn("if this message in docker conteiner, just skip it. But if it in local machine create .env file")
	}
}

func GetString(key string) string {
	return viper.GetString(key)
}
