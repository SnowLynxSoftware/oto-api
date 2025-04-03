package cmd

import (
	"github.com/snowlynxsoftware/oto-api/config"
	"github.com/snowlynxsoftware/oto-api/server"
)

type ServerCommand struct {
}

func (s *ServerCommand) Execute() error {
	appConfig := config.NewAppConfig()

	server := server.NewAppServer(appConfig)
	server.Start()
	return nil
}
