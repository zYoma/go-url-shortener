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

const (
	envServerAddress = "SERVER_ADDRESS"
	envBaseURL       = "BASE_URL"
	envLoggerLevel   = "LOG_LEVEL"
	envStorageFile   = "FILE_STORAGE_PATH"
	envDSN           = "DATABASE_DSN"
	envTokenSecret   = "TOKEN_SECRET"
)

type Config struct {
	RunAddr      string
	BaseShortURL string
	LogLevel     string
	StorageFile  string
	DSN          string
	TokenSecret  string
}

func GetConfig() *Config {
	// парсим аргументы командной строки
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&flagBaseShortURL, "b", "http://localhost:8080", "base short url")
	flag.StringVar(&flagLogLevel, "l", "info", "log level")
	flag.StringVar(&flagStorageFileNmae, "f", "/tmp/short-url-db.json", "starage file name")
	flag.StringVar(&flagDSN, "d", "", "DB DSN")
	flag.StringVar(&flagTokenSecret, "s", "secret_for_test_only", "secret for jwt")
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

	return &Config{
		RunAddr:      flagRunAddr,
		BaseShortURL: flagBaseShortURL,
		LogLevel:     flagLogLevel,
		StorageFile:  flagStorageFileNmae,
		DSN:          flagDSN,
		TokenSecret:  flagTokenSecret,
	}
}
