package config

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type IAppConfig interface {
	GetCloudEnv() string
	IsDebugMode() bool
	GetDBConnectionString() string
	GetAuthHashPepper() string
	GetJWTSecretKey() string
	GetSendgridAPIKey() string
}

type AppConfig struct {
	cloudEnv           string
	debugMode          bool
	dBConnectionString string
	authHashPepper     string
	jwtSecretKey       string
	sendgridAPIKey     string
}

func NewAppConfig() IAppConfig {

	appConfig := &AppConfig{}
	// Required environment variables
	appConfig.cloudEnv = os.Getenv("CLOUD_ENV")

	// Default values
	appConfig.debugMode = os.Getenv("DEBUG_MODE") == "true"
	appConfig.dBConnectionString = ""
	appConfig.authHashPepper = ""
	appConfig.jwtSecretKey = ""
	appConfig.sendgridAPIKey = ""

	if appConfig.cloudEnv == "" {
		log.Fatal("[CLOUD_ENV] is required")
	}

	// Load any additional variables from the environment and override the secret manager values
	appConfig.dBConnectionString = os.Getenv("DB_CONNECTION_STRING")
	appConfig.authHashPepper = os.Getenv("AUTH_HASH_PEPPER")
	appConfig.jwtSecretKey = os.Getenv("JWT_SECRET_KEY")
	appConfig.sendgridAPIKey = os.Getenv("SENDGRID_API_KEY")

	errorList := ""

	if appConfig.dBConnectionString == "" {
		errorList += "[DB_CONNECTION_STRING]\n"
	}

	if appConfig.authHashPepper == "" {
		errorList += "[AUTH_HASH_PEPPER]\n"
	}

	if appConfig.jwtSecretKey == "" {
		errorList += "[JWT_SECRET_KEY]\n"
	}

	if appConfig.sendgridAPIKey == "" {
		errorList += "[SENDGRID_API_KEY]\n"
	}

	if errorList != "" {
		errorList = "Missing environment variables:\n" + errorList
		panic(errorList)
	}

	return appConfig
}

func (a *AppConfig) GetCloudEnv() string {
	return a.cloudEnv
}

func (a *AppConfig) IsDebugMode() bool {
	return a.debugMode
}

func (a *AppConfig) GetDBConnectionString() string {
	return a.dBConnectionString
}

func (a *AppConfig) GetAuthHashPepper() string {
	return a.authHashPepper
}

func (a *AppConfig) GetJWTSecretKey() string {
	return a.jwtSecretKey
}

func (a *AppConfig) GetSendgridAPIKey() string {
	return a.sendgridAPIKey
}
