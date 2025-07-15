package esp32

import (
	"encoding/json"
	"fmt"
	"log"
	"smart-fridge/internal/services/audio"
	"smart-fridge/internal/services/openai"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	lightSensorThresholdHigh = 3500
	lightSensorThresholdLow  = 1000
	greetingCooldown         = 15 * time.Second
)

// ESP32Status represents sensor data from the ESP32 device
type ESP32Status struct {
	Timestamp      int64 `json:"timestamp"`
	IsPeriodic     bool  `json:"isPeriodic"`
	LightSensor    int   `json:"lightSensor"`
	Button1Pressed bool  `json:"button_1_pressed"`
	Button2Pressed bool  `json:"button_2_pressed"`
	IsPlayingAudio bool  `json:"isPlayingAudio"`
}

type Dependencies struct {
	AudioService *audio.AudioService
	OpenAIClient *openai.OpenAIClient
	NatsPort     int
}

type DoorState int

const (
	DoorClosed DoorState = iota
	DoorOpening
	DoorOpen
)

type ESP32Service struct {
	Deps             Dependencies
	nats             *nats.Conn
	mutex            sync.Mutex
	currentDoorState DoorState
	lastStatus       *ESP32Status
	lastGreetingTime time.Time
	audioQueue       []string
	isAudioPlaying   bool
}

func NewESP32Service(deps Dependencies) *ESP32Service {
	return &ESP32Service{
		Deps: deps,
		nats: nil,
	}
}

func (s *ESP32Service) Start() error {
	nc, err := nats.Connect(fmt.Sprintf("127.0.0.1:%d", s.Deps.NatsPort))
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}
	s.nats = nc

	s.nats.Subscribe("esp32.sensors", func(msg *nats.Msg) {
		s.handleStatusUpdate(msg)
	})

	s.nats.Subscribe("esp32.image", func(msg *nats.Msg) {
		s.handleImageUpdate(msg)
	})

	return nil
}

func (s *ESP32Service) Stop() error {
	if s.nats != nil {
		s.nats.Close()
		s.nats = nil
	}
	return nil
}

func (s *ESP32Service) QueueAudio(path string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.audioQueue = append(s.audioQueue, path)

	log.Printf("Audio queued: %s", path)
	s.processAudioQueue()
}

func (s *ESP32Service) processAudioQueue() {
	if s.isAudioPlaying || len(s.audioQueue) == 0 {
		return
	}

	nextAudio := s.audioQueue[0]
	s.audioQueue = s.audioQueue[1:]

	if s.nats != nil {
		s.nats.Publish("esp32.play_audio", []byte(nextAudio))
		s.isAudioPlaying = true
		log.Printf("Sent audio to ESP32: %s", nextAudio)
	}
}

func (s *ESP32Service) updateDoorState(lightSensor int) DoorState {
	switch s.currentDoorState {
	case DoorClosed:
		if lightSensor > lightSensorThresholdHigh {
			return DoorOpening
		}
	case DoorOpening:
		return DoorOpen
	case DoorOpen:
		if lightSensor < lightSensorThresholdLow {
			return DoorClosed
		}
	}
	return s.currentDoorState
}

func (s *ESP32Service) shouldPlayGreeting(newState DoorState) bool {
	if s.currentDoorState == DoorClosed && newState == DoorOpening {
		if time.Since(s.lastGreetingTime) > greetingCooldown {
			return true
		}
	}
	return false
}

// playRandomGreeting queues a random greeting from the audio service
func (s *ESP32Service) playRandomGreeting() {
	s.QueueAudio("/audio/persitent/greeting")
}

func (s *ESP32Service) handleStatusUpdate(msg *nats.Msg) {
	var status ESP32Status
	if err := json.Unmarshal(msg.Data, &status); err != nil {
		log.Printf("Failed to unmarshal ESP32 status: %v", err)
		return
	}

	s.mutex.Lock()
	s.lastStatus = &status

	// Update audio playing status
	wasPlaying := s.isAudioPlaying
	s.isAudioPlaying = status.IsPlayingAudio

	// If audio just finished playing, try to play next in queue
	if wasPlaying && !s.isAudioPlaying {
		log.Println("Audio playback finished, processing queue")
		s.processAudioQueue()
	}
	s.mutex.Unlock()

	// Skip periodic status updates for now
	if status.IsPeriodic {
		return
	}

	// All non periodic updates are processed here
	newDoorState := s.updateDoorState(status.LightSensor)

	// Check if we should play a greeting
	if s.shouldPlayGreeting(newDoorState) {
		log.Printf("Door opening detected (light: %d), playing greeting", status.LightSensor)
		s.playRandomGreeting()
	}

	// Update current state
	if newDoorState != s.currentDoorState {
		log.Printf("Door state changed: %s -> %s", s.currentDoorState, newDoorState)
		s.currentDoorState = newDoorState
	}

	// Handle button presses
	if status.Button1Pressed {
		log.Println("Button 1 pressed")
		s.handleButton1Press()
	}

	if status.Button2Pressed {
		log.Println("Button 2 pressed")
		s.handleButton2Press()
	}

	log.Printf("ESP32 Status - Light: %d, Door: %s, Audio: %t, Buttons: [%t,%t]",
		status.LightSensor, s.currentDoorState, status.IsPlayingAudio,
		status.Button1Pressed, status.Button2Pressed)
}

func (s *ESP32Service) handleImageUpdate(msg *nats.Msg) {
	log.Printf("Received image data from ESP32: %d bytes", len(msg.Data))
}

func (s *ESP32Service) handleButton1Press() {
}

func (s *ESP32Service) handleButton2Press() {

}

func (s *ESP32Service) GetStatus() *ESP32Status {
	if s.lastStatus == nil {
		return &ESP32Status{}
	}
	return s.lastStatus
}
