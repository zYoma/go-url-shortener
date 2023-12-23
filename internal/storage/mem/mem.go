package mem

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/zYoma/go-url-shortener/internal/config"
	"github.com/zYoma/go-url-shortener/internal/logger"
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

// Читает данные из файла при старте приложения
func (s *Storage) Init(cfg *config.Config) error {
	// открываем файл для чтения
	file, err := os.OpenFile(cfg.StorageFile, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Log.Sugar().Errorf("Не удалось открыть файл файл: %s", err)
		return err
	}
	defer file.Close()

	// Получаем информацию о файле для проверки размера
	fileInfo, err := file.Stat()
	if err != nil {
		logger.Log.Sugar().Errorf("Ошибка получения информации о файле: %s", err)
		return err
	}

	// Сброс указателя чтения файла в начало
	// if _, err := file.Seek(0, 0); err != nil {
	// 	logger.Log.Sugar().Errorf("Ошибка сброса указателя чтения файла: %s", err)
	// 	return err
	// }
	// Декодируем содержимое файла в map[string]string

	if fileInfo.Size() > 0 {
		if err := json.NewDecoder(file).Decode(&s.db); err != nil {
			logger.Log.Sugar().Errorf("Ошибка декодирования JSON: %s", err)
			return err
		}
	}

	return nil
}

// Пишет данные в файл при остановки приложения
func (s *Storage) Stop(cfg *config.Config) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.db) > 0 {
		// открываем файл для записи
		file, err := os.OpenFile(cfg.StorageFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			logger.Log.Sugar().Errorf("Не удалось открыть файл файл: %s", err)
			return err
		}
		defer file.Close()

		// Сериализуем map в JSON и записываем в файл
		if err := json.NewEncoder(file).Encode(&s.db); err != nil {
			logger.Log.Sugar().Errorf("Ошибка записи в файл: %s", err)
			return err
		}

		logger.Log.Sugar().Info("Данные успешно записаны в файл.")
	}

	return nil
}
