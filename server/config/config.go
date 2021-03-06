package config

import "github.com/spf13/viper"

type AppConfig struct {
	HTTPPort         int
	MySQLConfig      MySQLConfig
	NotiProdiverType string
}

func Load() *AppConfig {
	v := viper.New()
	v.SetConfigType("json")
	v.AutomaticEnv()

	mySQLConfig := MySQLConfig{
		Host:     v.GetString("MYSQL_HOST"),
		Username: v.GetString("MYSQL_USERNAME"),
		Password: v.GetString("MYSQL_PASSWORD"),
		Database: v.GetString("MYSQL_DATABASE"),
		Port:     v.GetInt("MYSQL_PORT"),
		Options:  "?charset=utf8&parseTime=True",
	}
	v.SetDefault("HTTP_PORT", 8080)
	v.SetDefault("NOTI_PROVIDER_TYPE", "email")

	return &AppConfig{
		MySQLConfig:      mySQLConfig,
		HTTPPort:         v.GetInt("HTTP_PORT"),
		NotiProdiverType: v.GetString("NOTI_PROVIDER_TYPE"),
	}
}
