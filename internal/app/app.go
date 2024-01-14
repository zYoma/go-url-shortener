package app

import (
	"github.com/zYoma/go-url-shortener/internal/app/server"
	"github.com/zYoma/go-url-shortener/internal/config"
	"github.com/zYoma/go-url-shortener/internal/logger"
	"github.com/zYoma/go-url-shortener/internal/storage"
	"github.com/zYoma/go-url-shortener/internal/storage/mem"
	"github.com/zYoma/go-url-shortener/internal/storage/postgres"
)

type App struct {
	Server *server.HTTPServer
}

func New(cfg *config.Config) (*App, error) {
	// создаем провайдер для storage
	provider, err := StorageConstructor(cfg)
	if err != nil {
		return nil, err
	}

	// инициализируем провайдера
	if err := provider.Init(); err != nil {
		return nil, err
	}

	// создаем сервер
	server := server.New(provider, cfg)

	return &App{Server: server}, nil
}

func (s *App) Run() error {
	// Создание канала для ошибок
	errChan := make(chan error)

	// запустить сервис
	logger.Log.Info("start application")

	go func() {
		errChan <- s.Server.Run()
	}()

	// Ожидание и обработка ошибки из горутины
	if err := <-errChan; err != nil {
		return err
	}

	return nil

}

// в зависимости от конфигурации, выбирает провайдера
func StorageConstructor(cfg *config.Config) (storage.StorageProvider, error) {
	if cfg.DSN != "" {
		logger.Log.Sugar().Infof("провайдер - postgres")
		return postgres.New(cfg)
	}
	logger.Log.Sugar().Infof("провайдер - mem")
	return mem.New(cfg)
}
