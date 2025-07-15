// handleWebUI serves the main web interface
package webui

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"smart-fridge/internal/services/audio"
	"smart-fridge/internal/services/esp32"
	"smart-fridge/internal/services/openai"

	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/starfederation/datastar-go/datastar"

	_ "embed"
)

type Dependencies struct {
	AudioService *audio.AudioService
	OpenAIClient *openai.OpenAIClient
	ESP32Service *esp32.ESP32Service
	NatsPort     int
}

type WebUIService struct {
	Deps Dependencies
	nats *nats.Conn
}

//go:embed ui.html
var uiHTML string

func NewWebUIService(deps Dependencies) *WebUIService {
	return &WebUIService{
		Deps: deps,
	}
}

func (s *WebUIService) Start() error {
	nc, err := nats.Connect(fmt.Sprintf("127.0.0.1:%d", s.Deps.NatsPort))
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}
	s.nats = nc

	return nil
}

func (s *WebUIService) RegisterWebUIRoutes(r *mux.Router) {
	// Audio service routes

	// WebUI routes
	r.HandleFunc("/", s.handleWebUI).Methods("GET")
	r.HandleFunc("/state", s.handleState).Methods("GET")
	r.HandleFunc("/led/", s.handleLed)
	r.HandleFunc("/play_tts", s.handleTTS)
	r.HandleFunc("/take_image", s.handleImage)
	r.HandleFunc("/greeting", s.handleGreeting)
}

func (s *WebUIService) handleWebUI(w http.ResponseWriter, r *http.Request) {
	// Serve the UI HTML
	w.Write([]byte(uiHTML))
}

// SSE handler: streams updates continuously
func (s *WebUIService) handleState(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)

	s.nats.Subscribe("esp32.image", func(msg *nats.Msg) {
		imageString := "data:image/jpg;base64," + string(msg.Data)
		sse.PatchSignals([]byte(fmt.Sprintf(`{"imageUrl": "%s"}`, imageString)))
	})

	s.nats.Subscribe("esp32.sensors", func(msg *nats.Msg) {
		sse.PatchSignals(msg.Data)
	})

	<-r.Context().Done()
}

func (s *WebUIService) handleLed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ledID := vars["id"]
	action := r.URL.Query().Get("action")

	if ledID == "" || action == "" {
		http.Error(w, "Invalid LED request", http.StatusBadRequest)
		return
	}

	topic := "esp32.led"
	command := fmt.Sprintf("%s %s", ledID, action)
	if err := s.nats.Publish(topic, []byte(command)); err != nil {
		http.Error(w, fmt.Sprintf("Failed to send LED command: %v", err),
			http.StatusInternalServerError)
		return
	}
	log.Printf("LED command sent: %s %s", ledID, action)

	w.WriteHeader(http.StatusOK)
}

type TTSRequest struct {
	MessageText string `json:"messageText"`
}

func (s *WebUIService) handleTTS(w http.ResponseWriter, r *http.Request) {

	var req TTSRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode request: %v", err), http.StatusBadRequest)
		return
	}

	audioData, err := s.Deps.OpenAIClient.TextToSpeech(req.MessageText)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate speech: %v", err), http.StatusInternalServerError)
		return
	}

	audioID, err := s.Deps.AudioService.StoreTemporaryAudio(audioData)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to store audio: %v", err), http.StatusInternalServerError)
		return
	}
	audioLink := fmt.Sprintf("/audio/%s", audioID)

	s.nats.Publish("esp32.play_audio", []byte(audioLink))

	w.WriteHeader(http.StatusOK)
}

func (s *WebUIService) handleImage(w http.ResponseWriter, r *http.Request) {
	s.nats.Publish("esp32.take_image", nil)
	w.WriteHeader(http.StatusOK)
	log.Println("Image capture request sent to ESP32")
}

func (s *WebUIService) handleGreeting(w http.ResponseWriter, r *http.Request) {
	s.nats.Publish("esp32.play_audio", []byte("/audio/persistent/greeting"))
	w.WriteHeader(http.StatusOK)
	log.Println("Greeting request sent to ESP32")
}
