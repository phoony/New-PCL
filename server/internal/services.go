package internal

import (
	"fmt"
	"log"
	"net/http"
	"smart-fridge/internal/services/audio"
	"smart-fridge/internal/services/esp32"
	"smart-fridge/internal/services/httpsrv"
	"smart-fridge/internal/services/nats"
	"smart-fridge/internal/services/openai"
	"smart-fridge/internal/services/webui"

	"github.com/gorilla/mux"
)

type Config struct {
	Addr         string
	NatsPort     int
	HTTPPort     int
	OpenAIAPIKey string
}

type Services struct {
	Config       Config
	AudioService *audio.AudioService
	OpenAIClient *openai.OpenAIClient
	NatsService  *nats.NATSService
	HTTPService  *httpsrv.HTTPService
	WebUIService *webui.WebUIService
	ESP32Service *esp32.ESP32Service

	router *mux.Router
}

func NewServices(config Config) *Services {
	return &Services{
		Config:       config,
		AudioService: audio.SetupAudioService(),
		OpenAIClient: openai.SetupOpenAIClient(config.OpenAIAPIKey),
		NatsService:  nats.NewNATSService(config.Addr, config.NatsPort),
		HTTPService: &httpsrv.HTTPService{
			Port: config.HTTPPort,
		},
		router: mux.NewRouter(),
	}
}

func (s *Services) Start() error {
	log.Println("Starting Smart Fridge services...")

	// Start NATS server
	if err := s.NatsService.Start(); err != nil {
		return fmt.Errorf("NATS startup failed: %w", err)
	}

	if err := s.HTTPService.StartHTTPServer(s.router); err != nil {
		return fmt.Errorf("HTTP server startup failed: %w", err)
	}

	// Start ESP32 service
	esp32Deps := esp32.Dependencies{
		AudioService: s.AudioService,
		OpenAIClient: s.OpenAIClient,
		NatsPort:     s.Config.NatsPort,
	}

	s.ESP32Service = esp32.NewESP32Service(esp32Deps)
	if err := s.ESP32Service.Start(); err != nil {
		return fmt.Errorf("ESP32 service startup failed: %w", err)
	}

	// Start Web UI service
	webuiDeps := webui.Dependencies{
		AudioService: s.AudioService,
		ESP32Service: s.ESP32Service,
		OpenAIClient: s.OpenAIClient,
		NatsPort:     s.Config.NatsPort,
	}
	s.WebUIService = webui.NewWebUIService(webuiDeps)
	if err := s.WebUIService.Start(); err != nil {
		return fmt.Errorf("Web UI service startup failed: %w", err)
	}

	// Setup HTTP server
	s.setupHTTPRoutes()
	s.AddLoggingMiddleware()

	log.Println("All services started successfully")
	return nil
}

func (s *Services) setupHTTPRoutes() {
	s.AudioService.RegisterAudioRoutes(s.router)
	s.WebUIService.RegisterWebUIRoutes(s.router)
}

func (s *Services) Stop() error {
	log.Println("Shutting down Smart Fridge services...")

	// Stop ESP32 service
	if err := s.ESP32Service.Stop(); err != nil {
		return fmt.Errorf("ESP32 service shutdown failed: %w", err)
	}

	// Shutdown HTTP server
	if err := s.HTTPService.Shutdown(); err != nil {
		return fmt.Errorf("HTTP server shutdown failed: %w", err)
	}

	// Shutdown NATS server
	s.NatsService.Stop()

	log.Println("All services stopped")
	return nil
}

func (s *Services) AddLoggingMiddleware() {
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
			next.ServeHTTP(w, r)
		})
	})
}
