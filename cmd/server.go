package cmd

import (
	"fmt"

	"github.com/sonwlynxsoftware/oto-api/config"
)

type ServerCommand struct {
}

func (s *ServerCommand) Execute() error {
	appConfig := config.NewAppConfig()

	fmt.Println("Hello, World! " + appConfig.GetCloudEnv() + " " + appConfig.GetDBConnectionString())

	return nil
}
