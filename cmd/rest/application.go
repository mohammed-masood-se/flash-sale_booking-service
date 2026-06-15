package main

import (
	"booking-service/internal/adapters/rest"
	"fmt"
	"log"
)

type ApplicationConfig struct {
	RestServerPort string
}

type Application struct {
	config *ApplicationConfig

	RestServer *rest.RestServer
}

func NewApplication(config ApplicationConfig) (*Application, error) {

	restServer, err := rest.NewRestServer(rest.RestServerConfig{
		Port: config.RestServerPort,
	})

	if err != nil {
		return nil, fmt.Errorf("failed creating NewRestServer: %w", err)
	}

	return &Application{
		config: &config,

		RestServer: restServer,
	}, nil
}

func (app *Application) Run() error {
	log.Printf("[rest-server] started on port %v\n", app.config.RestServerPort)
	err := app.RestServer.Run()

	return fmt.Errorf("failed running RestServer: %w", err)
}
