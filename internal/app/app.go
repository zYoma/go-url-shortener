package app

import (
	"github.com/zYoma/go-url-shortener/internal/app/server"
	"github.com/zYoma/go-url-shortener/internal/config"
	"github.com/zYoma/go-url-shortener/internal/handlers"
	"github.com/zYoma/go-url-shortener/internal/logger"
	"github.com/zYoma/go-url-shortener/internal/storage/mem"
)

type App struct {
	Server   *server.HTTPServer
	Provider handlers.URLProvider
}

func New(cfg *config.Config) (*App, error) {
	// создаем провайдер для storage
	provider := mem.New()

	// инициализируем провайдера
	err := provider.Init(cfg)
	if err != nil {
		return nil, err
	}

	// создаем сервер
	server := server.New(provider, cfg)

	return &App{Server: server, Provider: provider}, nil
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
