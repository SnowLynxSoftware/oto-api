package cmd

import (
	"github.com/sonwlynxsoftware/oto-api/config"
	"github.com/sonwlynxsoftware/oto-api/server"
)

type ServerCommand struct {
}

func (s *ServerCommand) Execute() error {
	appConfig := config.NewAppConfig()

	server := server.NewAppServer(appConfig)
	server.Start()
	return nil
}
