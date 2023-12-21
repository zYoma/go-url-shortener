package config

import (
	"flag"
	"os"
)

var flagRunAddr string
var flagBaseShortURL string
var flagLogLevel string

const (
	envServerAddress = "SERVER_ADDRESS"
	envBaseURL       = "BASE_URL"
	envLoggerLevel   = "LOG_LEVEL"
)

type Config struct {
	RunAddr      string
	BaseShortURL string
	LogLevel     string
}

func GetConfig() *Config {
	// парсим аргументы командной строки
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&flagBaseShortURL, "b", "http://localhost:8080", "base short url")
	flag.StringVar(&flagLogLevel, "l", "info", "log level")
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

	return &Config{
		RunAddr:      flagRunAddr,
		BaseShortURL: flagBaseShortURL,
		LogLevel:     flagLogLevel,
	}
}
