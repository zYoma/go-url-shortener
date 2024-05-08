package storage

import (
	"context"

	"github.com/zYoma/go-url-shortener/internal/models"
)

// StorageProvider определяет интерфейс для компонентов, отвечающих за хранение данных.
// Этот интерфейс включает в себя все методы URLProvider, что позволяет его реализациям
// обрабатывать операции с URL.
type StorageProvider interface {
	URLProvider
}

// URLProvider определяет набор методов для управления URL в хранилище, включая
// сохранение, извлечение и удаление URL, а также операции для работы с пакетами URL.
// Этот интерфейс предназначен для взаимодействия с различными реализациями хранилищ,
// поддерживающих операции с короткими и полными URL.
type URLProvider interface {
	// SaveURL сохраняет короткий и полный URL, ассоциированные с идентификатором пользователя.
	SaveURL(ctx context.Context, fullURL, shortURL, userID string) error

	// BulkSaveURL выполняет массовое сохранение данных о URL для указанного пользователя.
	BulkSaveURL(ctx context.Context, data []models.InsertData, userID string) error

	// GetURL извлекает полный URL по его короткой версии.
	GetURL(ctx context.Context, shortURL string) (string, error)

	// GetShortURL извлекает короткий URL по его полной версии.
	GetShortURL(ctx context.Context, shortURL string) (string, error)

	// Init инициализирует хранилище, подготавливая его к работе.
	Init() error

	// Ping проверяет доступность и работоспособность хранилища.
	Ping(ctx context.Context) error

	// GetUserURLs извлекает список всех URL, созданных пользователем.
	GetUserURLs(ctx context.Context, baseURL, userID string) ([]models.UserURLS, error)

	// DeleteListURL удаляет список URL, ассоциированных с идентификаторами пользователей.
	DeleteListURL(ctx context.Context, messages []models.UserListURLForDelete) error

	// GetServiceStats получает статистику сервиса.
	GetServiceStats(ctx context.Context) (models.ServiceStat, error)
}
