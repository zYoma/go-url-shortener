package postgres

import (
	"context"
	"database/sql"
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

// возможные ошибки пакета
var (
	// ErrCreatePool описывает ошибку создания пула соединений с базой данных.
	ErrCreatePool = errors.New("unable to create connection pool")
	// ErrPing описывает ошибку проверки соединения с базой данных.
	ErrPing = errors.New("checking connection to the database")
	// ErrURLNotFound описывает ошибку, возникающую, когда URL не найден в базе данных.
	ErrURLNotFound = errors.New("url not found")
	// ErrSaveURL описывает ошибку сохранения URL в базе данных.
	ErrSaveURL = errors.New("saving to database")
	// ErrCreateTable описывает ошибку создания таблиц в базе данных.
	ErrCreateTable = errors.New("creating tables")
	// ErrConflict описывает ошибку конфликта при попытке вставки URL, который уже существует.
	ErrConflict = errors.New("url already exist")
	// ErrGetURL описывает ошибку получения данных из базы данных.
	ErrGetURL = errors.New("select from database")
	// ErrScanRows описывает ошибку чтения строк из результата запроса.
	ErrScanRows = errors.New("scan rows")
	// ErrSRows описывает ошибку поиска строки в базе данных.
	ErrSRows = errors.New("line search error")
	// ErrUpdateURL описывает ошибку обновления данных о URL в базе данных.
	ErrUpdateURL = errors.New("update urls")
	// ErrURLDeleted описывает ошибку, возникающую при попытке доступа к удалённому URL.
	ErrURLDeleted = errors.New("URL was deleted")
)

// Storage реализует интерфейс StorageProvider и предоставляет методы для работы с хранилищем URL.
type Storage struct {
	pool  *pgxpool.Pool // Пул соединений с базой данных.
	mutex sync.Mutex    // Мьютекс для синхронизации доступа к базе данных.
}

// New инициализирует новый экземпляр Storage с подключением к базе данных, указанной в конфигурации.
func New(cfg *config.Config) (storage.StorageProvider, error) {
	dbpool, err := pgxpool.New(context.Background(), cfg.DSN)
	if err != nil {
		return nil, ErrCreatePool
	}
	return &Storage{pool: dbpool}, nil
}

// SaveURL сохраняет указанный URL в базе данных, ассоциируя его с конкретным пользователем.
func (s *Storage) SaveURL(ctx context.Context, fullURL string, shortURL string, userID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, err := s.pool.Exec(ctx, `
        INSERT INTO url (full_url, short_url, user_id) VALUES ($1, $2, $3) ;
    `, fullURL, shortURL, userID)

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

// GetShortURL возвращает короткий URL по заданному полному URL.
func (s *Storage) GetShortURL(ctx context.Context, fullURL string) (string, error) {
	var shortURL string
	row := s.pool.QueryRow(ctx, `SELECT short_url FROM url WHERE full_url = $1`, fullURL)
	err := row.Scan(&shortURL)
	if err != nil {
		return "", ErrURLNotFound
	}

	return shortURL, nil
}

// GetURL возвращает полный URL по заданному короткому URL.
func (s *Storage) GetURL(ctx context.Context, shortURL string) (string, error) {
	var (
		fullURL   string
		isDeleted bool
	)
	row := s.pool.QueryRow(ctx, `SELECT full_url, is_deleted FROM url WHERE short_url = $1`, shortURL)
	err := row.Scan(&fullURL, &isDeleted)
	if err != nil {
		// Если URL не найден, возвращаем соответствующую ошибку
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrURLNotFound
		}
		// Для других ошибок возвращаем их напрямую
		return "", err
	}

	// Проверяем, помечен ли URL как удаленный
	if isDeleted {
		return "", ErrURLDeleted
	}

	return fullURL, nil
}

// Init выполняет инициализацию хранилища, включая создание необходимых таблиц.
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
			"created" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			"user_id" UUID NOT NULL,
			"is_deleted" BOOLEAN DEFAULT FALSE
		);
    `)

	tx.Exec(ctx, `CREATE UNIQUE INDEX IF NOT EXISTS idx_full_url_unique ON url(full_url);`)

	return tx.Commit(ctx)
}

// Ping проверяет состояние соединения с базой данных.
func (s *Storage) Ping(ctx context.Context) error {
	if err := s.pool.Ping(ctx); err != nil {
		return ErrPing
	}
	return nil
}

// BulkSaveURL выполняет массовое сохранение данных о URL для указанного пользователя.
func (s *Storage) BulkSaveURL(ctx context.Context, data []models.InsertData, userID string) error {
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
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		valueArgs = append(valueArgs, d.OriginalURL, d.ShortURL, userID)
	}

	// Формирование и выполнение запроса
	stmt := fmt.Sprintf("INSERT INTO url (full_url, short_url, user_id) VALUES %s", strings.Join(valueStrings, ","))
	_, err := s.pool.Exec(ctx, stmt, valueArgs...)
	if err != nil {
		logger.Log.Sugar().Errorf("Не удалось сохранить url: %s", err)
		return ErrSaveURL
	}

	return nil
}

// GetUserURLs возвращает список URL, созданных пользователем.
func (s *Storage) GetUserURLs(ctx context.Context, baseURL string, userID string) ([]models.UserURLS, error) {
	var urls []models.UserURLS
	rows, err := s.pool.Query(ctx, `SELECT short_url, full_url FROM url WHERE user_id = $1`, userID)
	if err != nil {
		logger.Log.Sugar().Errorf("Не удалось выполнить запрос: %s", err)
		return nil, ErrGetURL
	}
	defer rows.Close()

	for rows.Next() {
		var pair models.UserURLS
		if err = rows.Scan(&pair.ShortURL, &pair.OriginalURL); err != nil {
			logger.Log.Sugar().Errorf("Не удалось прочитать строку: %s", err)
			return nil, ErrScanRows
		}
		pair.ShortURL = fmt.Sprintf("%s/%s", baseURL, pair.ShortURL)
		urls = append(urls, pair)
	}

	// Проверяем наличие ошибок после завершения перебора
	if err = rows.Err(); err != nil {
		logger.Log.Sugar().Errorf("Ошибка: %s", err)
		return nil, ErrSRows
	}

	return urls, nil
}

// DeleteListURL удаляет список URL для заданных пользователей.
func (s *Storage) DeleteListURL(ctx context.Context, messages []models.UserListURLForDelete) error {
	if len(messages) == 0 {
		// Нет данных для обработки
		return nil
	}

	var (
		placeholders []string
		args         []interface{}
		argCounter   = 1
	)

	for _, message := range messages {
		for _, url := range message.URLS {
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", argCounter, argCounter+1))
			args = append(args, url, message.UserID)
			argCounter += 2
		}
	}

	if len(placeholders) == 0 {
		return nil // Нет URL для обновления
	}

	query := fmt.Sprintf(`UPDATE url SET is_deleted = true WHERE (short_url, user_id) IN (%s)`,
		strings.Join(placeholders, ", "))

	_, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		logger.Log.Sugar().Errorf("Не удалось выполнить обновление: %s", err)
		return ErrUpdateURL
	}

	return nil
}
