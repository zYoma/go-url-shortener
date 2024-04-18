package config

import (
	"encoding/json"
	"flag"
	"os"
	"reflect"
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
var flagConfigFile string

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
	envConfigFile    = "CONFIG"
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

type fileConfig struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	LogLevel        string `json:"log_level"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	TokenSecret     string `json:"token_secret"`
	EnableHTTPS     bool   `json:"enable_https"`
	CertPath        string `json:"cert_path"`
	CertKeyPath     string `json:"cert_key_path"`
}

func parseConfigFile(filePath string) (*fileConfig, error) {
	if filePath == "" {
		return nil, nil
	}
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var conf fileConfig
	if err := json.Unmarshal(fileBytes, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

// GetConfig парсит аргументы командной строки и переменные окружения,
// создавая и возвращая конфигурацию приложения. Приоритет имеют значения из переменных окружения.
//
// Возвращает сконфигурированный экземпляр *Config.
func GetConfig() (*Config, error) {

	// парсим аргументы командной строки
	flag.StringVar(&flagRunAddr, "a", "", "address and port to run server")
	flag.StringVar(&flagBaseShortURL, "b", "", "base short url")
	flag.StringVar(&flagLogLevel, "l", "", "log level")
	flag.StringVar(&flagStorageFileNmae, "f", "", "starage file name")
	flag.StringVar(&flagDSN, "d", "", "DB DSN")
	flag.StringVar(&flagTokenSecret, "j", "", "secret for jwt")
	flag.BoolVar(&flagHTTPS, "s", false, "enable HTTPS")
	flag.StringVar(&flagCertPath, "cr", "", "path to cert")
	flag.StringVar(&flagCertKeyPath, "ck", "", "path to cert key")
	flag.StringVar(&flagConfigFile, "c", "", "path to config file")
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

	confFromFile, err := parseConfigFile(flagConfigFile)
	if err != nil {
		return nil, err
	}

	if confFromFile != nil {
		setValueFromFileConfig(&flagRunAddr, confFromFile.ServerAddress)
		setValueFromFileConfig(&flagBaseShortURL, confFromFile.BaseURL)
		setValueFromFileConfig(&flagLogLevel, confFromFile.LogLevel)
		setValueFromFileConfig(&flagStorageFileNmae, confFromFile.FileStoragePath)
		setValueFromFileConfig(&flagDSN, confFromFile.DatabaseDSN)
		setValueFromFileConfig(&flagTokenSecret, confFromFile.TokenSecret)
		setValueFromFileConfig(&flagHTTPS, confFromFile.EnableHTTPS)
		setValueFromFileConfig(&flagCertPath, confFromFile.CertPath)
		setValueFromFileConfig(&flagCertKeyPath, confFromFile.CertKeyPath)
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
	}, nil
}

func nilValue[T comparable]() T {
	var zero T
	return zero
}

// setValueFromFileConfig проставляет значения из файла конфигурации, если текущее значение пустое
// дженерики использовал чтобы работать как с bool так и со строкой
func setValueFromFileConfig[T comparable](varPtr *T, varFile T) {
	switch reflect.TypeOf(*varPtr).Kind() {
	case reflect.String:
		if *varPtr == nilValue[T]() && varFile != nilValue[T]() {
			*varPtr = varFile
		}
	case reflect.Bool:
		if *varPtr == nilValue[T]() {
			*varPtr = varFile
		}
	}
}
