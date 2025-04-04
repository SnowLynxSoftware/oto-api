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
}

type AppConfig struct {
	cloudEnv           string
	awsAccessKeyId     string
	awsSecretAccessKey string
	debugMode          bool
	dBConnectionString string
}

func NewAppConfig() IAppConfig {

	appConfig := &AppConfig{}
	// Required environment variables
	appConfig.cloudEnv = os.Getenv("CLOUD_ENV")
	appConfig.awsAccessKeyId = os.Getenv("AWS_ACCESS_KEY_ID")
	appConfig.awsSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

	// Default values
	appConfig.debugMode = os.Getenv("DEBUG_MODE") == "true"
	appConfig.dBConnectionString = ""

	if appConfig.cloudEnv == "" {
		log.Fatal("[CLOUD_ENV] is required")
	}

	useSecretManager := os.Getenv("USE_SECRET_MANAGER") == "true"

	// If we use secret manager, it will happen first.
	if useSecretManager {
		secretManager := NewSecretManagerConfig(appConfig.cloudEnv)
		appConfig.debugMode, _ = secretManager.GetDebugMode()
		appConfig.dBConnectionString, _ = secretManager.GetDBConnectionString()
	}

	// Load any additional variables from the environment and override the secret manager values
	appConfig.dBConnectionString = os.Getenv("DB_CONNECTION_STRING")

	errorList := ""

	if appConfig.dBConnectionString == "" {
		errorList += "[DB_CONNECTION_STRING]\n"
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
