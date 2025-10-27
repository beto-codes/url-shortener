package storage

import (
	"sync"

	"github.com/beto-codes/url-shortener/internal/utils/constants"
)

type MemoryStorage struct {
	mu   sync.RWMutex
	urls map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		urls: make(map[string]string),
	}
}

func (m *MemoryStorage) Save(shortCode, longURL string) error {
	if longURL == constants.EmptyString {
		return ErrEmpty
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.urls[shortCode] = longURL
	return nil
}

func (m *MemoryStorage) Get(shortCode string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	longURL, exists := m.urls[shortCode]
	if !exists {
		return constants.EmptyString, ErrNotFound
	}

	return longURL, nil
}

func (m *MemoryStorage) Exists(shortCode string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.urls[shortCode]
	return exists
}
