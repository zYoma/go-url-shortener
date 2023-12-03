package mem

import (
	"errors"
	"sync"
)

// реализация хранилища в памяти
type Storage struct {
	db    map[string]string
	mutex sync.Mutex
}

func New() *Storage {
	db := make(map[string]string)
	return &Storage{db: db}
}

// SaveUrl to db.
func (s *Storage) SaveURL(fullURL string, shortURL string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.db[shortURL] = fullURL
}

// GetUrl from db.
func (s *Storage) GetURL(shortURL string) (string, error) {
	fullURL, ok := s.db[shortURL]
	if !ok {
		return "", errors.New("url not found")
	}

	return fullURL, nil
}
