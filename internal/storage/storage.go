package storage

import "errors"

var (
	ErrNotFound = errors.New("URL not found")
	ErrEmpty    = errors.New("URL cannot be empty")
)

type Storage interface {
	Save(shortCode, longURL string) error
	Get(shortCode string) (string, error)
	Exists(shortCode string) bool
}
