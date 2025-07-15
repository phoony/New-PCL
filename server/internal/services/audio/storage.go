package audio

import (
	"errors"
	"sync"
)

type AudioStorage interface {
	Save(id string, data []byte) error
	Get(id string) ([]byte, error)
	Delete(id string) error
	ListKeysWithPrefix(prefix string) []string
}

type MemoryStorage struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string][]byte),
	}
}

func (m *MemoryStorage) Save(id string, data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[id] = data
	return nil
}

func (m *MemoryStorage) Get(id string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	d, ok := m.data[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return d, nil
}

func (m *MemoryStorage) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, id)
	return nil
}

func (m *MemoryStorage) ListKeysWithPrefix(prefix string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	keys := []string{}
	for k := range m.data {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			keys = append(keys, k)
		}
	}
	return keys
}
