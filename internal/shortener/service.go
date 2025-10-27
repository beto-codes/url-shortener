package shortener

import (
	"github.com/beto-codes/url-shortener/internal/storage"
	"github.com/beto-codes/url-shortener/internal/utils/constants"
)

type Service struct {
	storage   storage.Storage
	generator *Generator
}

func NewService(store storage.Storage) *Service {
	return &Service{
		storage:   store,
		generator: NewGenerator(),
	}
}

func (s *Service) Shorten(longURL string) (string, error) {
	shortCode := s.generator.Generate(longURL)

	for s.storage.Exists(shortCode) {
		shortCode = s.generator.GenerateRandom()
	}

	err := s.storage.Save(shortCode, longURL)
	if err != nil {
		return constants.EmptyString, err
	}

	return shortCode, nil
}

func (s *Service) Resolve(shortCode string) (string, error) {
	return s.storage.Get(shortCode)
}
