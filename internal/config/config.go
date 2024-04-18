package config

import (
	"flag"
	"os"
)

var flagRunAddr string
var flagBaseShortURL string
var flagLogLevel string
var flagStorageFileNmae string
var flagDSN string
var flagTokenSecret string
var flagHTTPS bool
var flagCertPath string
var flagCertKeyPath string

const (
	envServerAddress = "SERVER_ADDRESS"
	envBaseURL       = "BASE_URL"
	envLoggerLevel   = "LOG_LEVEL"
	envStorageFile   = "FILE_STORAGE_PATH"
	envDSN           = "DATABASE_DSN"
	envTokenSecret   = "TOKEN_SECRET"
	envHTTPS         = "ENABLE_HTTPS"
	envCertPath      = "CERT_PATH"
	envCertKeyPath   = "CERT_KEY_PATH"
)

// Config определяет конфигурацию приложения, собираемую из аргументов командной строки и переменных окружения.
type Config struct {
	RunAddr      string // Адрес и порт для запуска сервера.
	BaseShortURL string // Базовый URL для коротких ссылок.
	LogLevel     string // Уровень логирования.
	StorageFile  string // Имя файла для хранения данных.
	DSN          string // Data Source Name для подключения к БД.
	TokenSecret  string // Секрет для подписи JWT токенов.
	EnableHTTPS  bool   // Включить HTTPS
	CertPath     string // путь до файла с сертификатом
	CertKeyPath  string // путь до ключа
}

// GetConfig парсит аргументы командной строки и переменные окружения,
// создавая и возвращая конфигурацию приложения. Приоритет имеют значения из переменных окружения.
//
// Возвращает сконфигурированный экземпляр *Config.
func GetConfig() *Config {
	// парсим аргументы командной строки
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&flagBaseShortURL, "b", "http://localhost:8080", "base short url")
	flag.StringVar(&flagLogLevel, "l", "info", "log level")
	flag.StringVar(&flagStorageFileNmae, "f", "/tmp/short-url-db.json", "starage file name")
	flag.StringVar(&flagDSN, "d", "", "DB DSN")
	flag.StringVar(&flagTokenSecret, "j", "secret_for_test_only", "secret for jwt")
	flag.BoolVar(&flagHTTPS, "s", false, "enable HTTPS")
	flag.StringVar(&flagCertPath, "cr", "certs/cert.pem", "path to cert")
	flag.StringVar(&flagCertKeyPath, "ck", "certs/key.pem", "path to cert key")
	flag.Parse()

	// если есть переменные окружения, используем их значения
	if envRunAddr := os.Getenv(envServerAddress); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	if envBaseShortURL := os.Getenv(envBaseURL); envBaseShortURL != "" {
		flagBaseShortURL = envBaseShortURL
	}
	if envLogLevel := os.Getenv(envLoggerLevel); envLogLevel != "" {
		flagLogLevel = envLogLevel
	}
	if envStorageFileName := os.Getenv(envStorageFile); envStorageFileName != "" {
		flagStorageFileNmae = envStorageFileName
	}
	if envDBDSN := os.Getenv(envDSN); envDBDSN != "" {
		flagDSN = envDBDSN
	}
	if envJWTSecret := os.Getenv(envTokenSecret); envJWTSecret != "" {
		flagTokenSecret = envJWTSecret
	}
	envEnableHTTPS := os.Getenv("ENABLE_HTTPS")
	if envEnableHTTPS != "" {
		flagHTTPS = (envEnableHTTPS == "1")
	}
	if envCert := os.Getenv(envCertPath); envCert != "" {
		flagCertPath = envCert
	}
	if envCertKey := os.Getenv(envCertKeyPath); envCertKey != "" {
		flagCertKeyPath = envCertKey
	}

	return &Config{
		RunAddr:      flagRunAddr,
		BaseShortURL: flagBaseShortURL,
		LogLevel:     flagLogLevel,
		StorageFile:  flagStorageFileNmae,
		DSN:          flagDSN,
		TokenSecret:  flagTokenSecret,
		EnableHTTPS:  flagHTTPS,
		CertPath:     flagCertPath,
		CertKeyPath:  flagCertKeyPath,
	}
}
