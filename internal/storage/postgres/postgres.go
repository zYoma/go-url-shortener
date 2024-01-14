package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
var ErrConflict = errors.New("url already exist")

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
        INSERT INTO url (full_url, short_url) VALUES ($1, $2) ;
    `, fullURL, shortURL)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return ErrConflict
		}
		logger.Log.Sugar().Errorf("Не удалось сохранить url: %s", err)
		return ErrSaveURL
	}
	return nil
}

func (s *Storage) GetShortURL(ctx context.Context, fullURL string) (string, error) {
	var shortURL string
	row := s.pool.QueryRow(ctx, `SELECT short_url FROM url WHERE full_url = $1`, fullURL)
	err := row.Scan(&shortURL)
	if err != nil {
		return "", ErrURLNotFound
	}

	return shortURL, nil
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

	tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS url (
			"id" SERIAL PRIMARY KEY,
			"full_url" VARCHAR(250) NOT NULL,
			"short_url" VARCHAR(250) NOT NULL,
			"created" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
    `)

	tx.Exec(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS idx_full_url_unique ON url(full_url);`)

	return tx.Commit(ctx)
}

func (s *Storage) Ping(ctx context.Context) error {
	if err := s.pool.Ping(ctx); err != nil {
		return ErrPing
	}
	return nil
}

func (s *Storage) BulkSaveURL(ctx context.Context, data []models.InsertData) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Проверка на пустой слайс
	if len(data) == 0 {
		return nil
	}

	// Начало подготовки запроса
	valueStrings := make([]string, 0, len(data))
	valueArgs := make([]interface{}, 0, len(data)*2)
	for i, d := range data {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		valueArgs = append(valueArgs, d.OriginalURL, d.ShortURL)
	}

	// Формирование и выполнение запроса
	stmt := fmt.Sprintf("INSERT INTO url (full_url, short_url) VALUES %s", strings.Join(valueStrings, ","))
	_, err := s.pool.Exec(ctx, stmt, valueArgs...)
	if err != nil {
		logger.Log.Sugar().Errorf("Не удалось сохранить url: %s", err)
		return ErrSaveURL
	}

	return nil
}
