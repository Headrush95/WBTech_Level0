package configs

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfig() error {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	if err := godotenv.Load(); err != nil {
		return err
	}
	return viper.ReadInConfig()
}
