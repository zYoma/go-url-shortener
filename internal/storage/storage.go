package storage

import (
	"context"
)

type StorageProvider interface {
	URLProvider
}

type URLProvider interface {
	SaveURL(ctx context.Context, fullURL string, shortURL string) error
	GetURL(ctx context.Context, shortURL string) (string, error)
	Init() error
	Ping(ctx context.Context) error
}
