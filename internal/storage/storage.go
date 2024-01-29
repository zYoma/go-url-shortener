package storage

import (
	"context"

	"github.com/zYoma/go-url-shortener/internal/models"
)

type StorageProvider interface {
	URLProvider
}

type URLProvider interface {
	SaveURL(ctx context.Context, fullURL string, shortURL string, userID string) error
	BulkSaveURL(ctx context.Context, data []models.InsertData, userID string) error
	GetURL(ctx context.Context, shortURL string) (string, error)
	GetShortURL(ctx context.Context, shortURL string) (string, error)
	Init() error
	Ping(ctx context.Context) error
	GetUserURLs(ctx context.Context, baseURL string, userID string) ([]models.UserURLS, error)
	DeleteListURL(ctx context.Context, messages []models.UserListURLForDelete) error
}
