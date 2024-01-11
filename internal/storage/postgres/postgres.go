package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zYoma/go-url-shortener/internal/config"
	"github.com/zYoma/go-url-shortener/internal/logger"
	"github.com/zYoma/go-url-shortener/internal/models"
	"github.com/zYoma/go-url-shortener/internal/storage"
)

var ErrCreatePool = errors.New("unable to create connection pool")
var ErrPing = errors.New("error when checking connection to the database")
var ErrURLNotFound = errors.New("url not found")
var ErrSaveURL = errors.New("error when saving to database")
var ErrCreateTable = errors.New("error creating tables")

type Storage struct {
	pool  *pgxpool.Pool
	mutex sync.Mutex
}

func New(cfg *config.Config) (storage.StorageProvider, error) {
	dbpool, err := pgxpool.New(context.Background(), cfg.DSN)
	if err != nil {
		return nil, ErrCreatePool
	}
	return &Storage{pool: dbpool}, nil
}

func (s *Storage) SaveURL(ctx context.Context, fullURL string, shortURL string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, err := s.pool.Exec(ctx, `
        INSERT INTO url (full_url, short_url) VALUES ($1, $2);
    `, fullURL, shortURL)

	if err != nil {
		logger.Log.Sugar().Errorf("Не удалось сохранить url: %s", err)
		return ErrSaveURL
	}
	return nil
}

func (s *Storage) GetURL(ctx context.Context, shortURL string) (string, error) {
	var fullURL string
	row := s.pool.QueryRow(ctx, `SELECT full_url FROM url WHERE short_url = $1`, shortURL)
	err := row.Scan(&fullURL)
	if err != nil {
		return "", ErrURLNotFound
	}

	return fullURL, nil
}

func (s *Storage) Init() error {
	ctx := context.Background()
	txOptions := pgx.TxOptions{}

	tx, err := s.pool.BeginTx(ctx, txOptions)
	if err != nil {
		logger.Log.Sugar().Errorf("Не удалось создать таблицу: %s", err)
		return ErrCreateTable
	}

	defer tx.Rollback(ctx)

	tx.Exec(context.TODO(), `
		CREATE TABLE IF NOT EXISTS url (
			"id" SERIAL PRIMARY KEY,
			"full_url" VARCHAR(250) NOT NULL,
			"short_url" VARCHAR(250) NOT NULL,
			"created" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
    `)

	return tx.Commit(ctx)
}

func (s *Storage) Ping(ctx context.Context) error {
	if err := s.pool.Ping(ctx); err != nil {
		return ErrPing
	}
	return nil
}

func (s *Storage) BulkSaveURL(ctx context.Context, data *[]models.InsertData) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	realData := *data

	// Построение строки SQL для вставки
	sqlValues := make([]string, len(realData))
	for i, d := range realData {
		sqlValues[i] = fmt.Sprintf("('%s', '%s')", d.OriginalURL, d.ShortURL)
	}

	query := fmt.Sprintf("INSERT INTO url (full_url, short_url) VALUES %s;", strings.Join(sqlValues, ", "))

	_, err := s.pool.Exec(ctx, query)
	if err != nil {
		logger.Log.Sugar().Errorf("Не удалось сохранить url: %s", err)
		return ErrSaveURL
	}
	return nil
}
