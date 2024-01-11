package storage

import (
	"context"

	"github.com/zYoma/go-url-shortener/internal/models"
)

type StorageProvider interface {
	URLProvider
}

type URLProvider interface {
	SaveURL(ctx context.Context, fullURL string, shortURL string) error
	BulkSaveURL(ctx context.Context, data *[]models.InsertData) error
	GetURL(ctx context.Context, shortURL string) (string, error)
	Init() error
	Ping(ctx context.Context) error
}
