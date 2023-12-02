package mem

import "errors"

type Storage struct {
	db map[string]string
}

func New() *Storage {
	db := make(map[string]string)
	return &Storage{db: db}
}

// SaveUrl to db.
func (s *Storage) SaveUrl(fullURL string, shortURL string) {
	s.db[shortURL] = fullURL
}

// GetUrl from db.
func (s *Storage) GetUrl(shortURL string) (string, error) {
	fullURL, ok := s.db[shortURL]
	if !ok {
		return "", errors.New("url not found")
	}

	return fullURL, nil
}
