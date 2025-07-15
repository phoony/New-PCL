package main

import (
	"log"
	"os"
	"os/signal"
	"smart-fridge/internal/application"
	"strconv"
	"syscall"
)

func main() {
	// Get ports from environment variables with defaults
	httpPort := 5000
	if envHTTPPort := os.Getenv("HTTP_PORT"); envHTTPPort != "" {
		if port, err := strconv.Atoi(envHTTPPort); err == nil {
			httpPort = port
		} else {
			log.Printf("Invalid HTTP_PORT value: %s, using default 5000", envHTTPPort)
		}
	}

	natsPort := 4222
	if envNATSPort := os.Getenv("NATS_PORT"); envNATSPort != "" {
		if port, err := strconv.Atoi(envNATSPort); err == nil {
			natsPort = port
		} else {
			log.Printf("Invalid NATS_PORT value: %s, using default 4222", envNATSPort)
		}
	}

	log.Printf("Starting Smart Fridge - HTTP: %d, NATS: %d", httpPort, natsPort)

	app := application.NewApplication(application.Config{
		Addr:         "0.0.0.0",
		NatsPort:     natsPort,
		HTTPPort:     httpPort,
		OpenAIAPIKey: os.Getenv("OPENAI_API_KEY"),
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	if err := app.Start(); err != nil {
		log.Fatalf("Error starting application: %v", err)
	}

	<-c
	signal.Stop(c)

	if err := app.Shutdown(); err != nil {
		log.Fatalf("Error shutting down application: %v", err)
	}

	os.Exit(0)
}
