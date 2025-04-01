package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type AppConfig struct {
	dBConnectionString string
}

func NewAppConfig() *AppConfig {

	var dbConnectionString =  os.Getenv("DB_CONNECTION_STRING");

	return &AppConfig{
		dBConnectionString: dbConnectionString,
	}
}
