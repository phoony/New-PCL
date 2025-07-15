package application

import (
	"fmt"
	"smart-fridge/internal"
)

type Config struct {
	Addr         string
	NatsPort     int
	HTTPPort     int
	OpenAIAPIKey string
}

type Application struct {
	config   Config
	services *internal.Services
}

func NewApplication(config Config) *Application {
	services := internal.NewServices(internal.Config{
		Addr:         config.Addr,
		NatsPort:     config.NatsPort,
		HTTPPort:     config.HTTPPort,
		OpenAIAPIKey: config.OpenAIAPIKey,
	})

	return &Application{
		config:   config,
		services: services,
	}
}

func (app *Application) Start() error {
	// Start all services
	if err := app.services.Start(); err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}

	return nil
}

func (app *Application) Shutdown() error {
	// Stop all services
	if err := app.services.Stop(); err != nil {
		return fmt.Errorf("failed to stop services: %w", err)
	}
	return nil
}
