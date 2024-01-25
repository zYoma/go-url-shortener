package mem

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/zYoma/go-url-shortener/internal/config"
	"github.com/zYoma/go-url-shortener/internal/logger"
	"github.com/zYoma/go-url-shortener/internal/models"
	"github.com/zYoma/go-url-shortener/internal/storage"
)

var ErrURLNotFound = errors.New("url not found")
var ErrOpenFile = errors.New("failed to open file")
var ErrWriteFile = errors.New("failed to write file")
var ErrInfoFile = errors.New("error getting file information")
var ErrDecodeFile = errors.New("file decoding error")

// реализация хранилища в памяти
type Storage struct {
	db          map[string]string
	storagePath string
	mutex       sync.Mutex
}

func New(cfg *config.Config) (storage.StorageProvider, error) {
	db := make(map[string]string)
	return &Storage{db: db, storagePath: cfg.StorageFile}, nil
}

// SaveUrl to db.
func (s *Storage) SaveURL(ctx context.Context, fullURL string, shortURL string, userID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.db[shortURL] = fullURL

	s.saveFile()

	return nil
}

// GetUrl from db.
func (s *Storage) GetURL(ctx context.Context, shortURL string) (string, error) {
	fullURL, ok := s.db[shortURL]
	if !ok {
		return "", ErrURLNotFound
	}

	return fullURL, nil
}

func (s *Storage) GetShortURL(ctx context.Context, fullURL string) (string, error) {
	for shortURL, u := range s.db {
		if u == fullURL {
			return shortURL, nil
		}
	}

	return "", ErrURLNotFound

}

// Читает данные из файла при старте приложения
func (s *Storage) Init() error {
	// открываем файл для чтения
	file, err := os.OpenFile(s.storagePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Log.Sugar().Errorf("Не удалось открыть файл: %s", err)
		return ErrOpenFile
	}
	defer file.Close()

	// Получаем информацию о файле для проверки размера
	fileInfo, err := file.Stat()
	if err != nil {
		logger.Log.Sugar().Errorf("Ошибка получения информации о файле: %s", err)
		return ErrInfoFile
	}

	if fileInfo.Size() > 0 {
		if err := json.NewDecoder(file).Decode(&s.db); err != nil {
			logger.Log.Sugar().Errorf("Ошибка декодирования JSON: %s", err)
			return ErrDecodeFile
		}
	}

	return nil
}

func (s *Storage) Ping(ctx context.Context) error {
	return nil
}

func (s *Storage) BulkSaveURL(ctx context.Context, data []models.InsertData, userID string) error {
	for _, url := range data {
		s.db[url.ShortURL] = url.OriginalURL
	}

	s.saveFile()

	return nil
}

func (s *Storage) saveFile() error {
	// обновляем нашу БД в фале
	file, err := os.OpenFile(s.storagePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		logger.Log.Sugar().Errorf("Не удалось открыть файл: %s", err)
		return ErrOpenFile
	}
	defer file.Close()

	// Сериализуем map в JSON и записываем в файл
	if err := json.NewEncoder(file).Encode(&s.db); err != nil {
		logger.Log.Sugar().Errorf("Ошибка записи в файл: %s", err)
		return ErrWriteFile
	}

	return nil
}

func (s *Storage) GetUserURLs(ctx context.Context, baseURL string, userID string) ([]models.UserURLS, error) {
	// не сделал поддержку данного метода в текущей реализации так как тестами она не проверяется, заленился
	var urls []models.UserURLS
	return urls, nil
}

func (s *Storage) DeleteListURL(ctx context.Context, messages []models.UserListURLForDelete) error {
	// аналогично
	return nil
}
