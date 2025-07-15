package audio

import (
	"bytes"
	crypto_rand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

type AudioService struct {
	storage AudioStorage
}

func NewAudioService(storage AudioStorage) *AudioService {
	return &AudioService{
		storage: storage,
	}
}

func generateID() string {
	b := make([]byte, 8)
	crypto_rand.Read(b)
	return hex.EncodeToString(b)
}

// Temporary/unique file
func (s *AudioService) StoreTemporaryAudio(data []byte) (id string, err error) {
	id = generateID()
	if err := s.storage.Save(id, data); err != nil {
		return "", fmt.Errorf("failed to store audio: %w", err)
	}
	return id, nil
}

// Persistent voiceline
func (s *AudioService) StoreNamedAudio(category, name string, data []byte) error {
	key := fmt.Sprintf("%s_%s", category, name)
	return s.storage.Save(key, data)
}

func (s *AudioService) DeleteAudio(id string) error {
	return s.storage.Delete(id)
}

func (s *AudioService) GetAudio(id string) ([]byte, error) {
	return s.storage.Get(id)
}

func (s *AudioService) GetRandomVoiceline(category string) ([]byte, string, error) {
	prefix := fmt.Sprintf("%s_", category)
	keys := s.storage.ListKeysWithPrefix(prefix)
	if len(keys) == 0 {
		return nil, "", errors.New("no voicelines found")
	}
	key := keys[rand.Intn(len(keys))]
	data, err := s.storage.Get(key)
	return data, key, err
}

func (s *AudioService) RegisterAudioRoutes(r *mux.Router) {
	r.HandleFunc("/audio/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		data, err := s.GetAudio(id)
		if err != nil {
			http.Error(w, "Audio not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "audio/mpeg")
		w.Header().Set("Content-Length", fmt.Sprint(len(data)))
		w.Header().Set("Accept-Ranges", "bytes")

		// Stream with notification on full read
		_, err = io.Copy(w, bytes.NewReader(data))
		if err != nil {
			log.Printf("Error streaming audio %s: %v", id, err)
		}
	})

	r.HandleFunc("/audio/persistent/{category}", func(w http.ResponseWriter, r *http.Request) {
		category := mux.Vars(r)["category"]

		data, _, err := s.GetRandomVoiceline(category)
		if err != nil {
			http.Error(w, "Audio not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "audio/mpeg")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
}
