package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

var ErrServerStoped = errors.New("server stoped")

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
	// Контекст с отменой для остановки сервера
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	// Создание канала для ошибок
	errChan := make(chan error)

	// запустить сервис
	logger.Log.Info("start application")
	go func() {
		if err := s.Server.Run(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Ожидание сигнала завершения работы
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		// При получении сигнала завершения останавливаем сервер
		if err := s.Server.Shutdown(ctx); err != nil {
			return err
		}
		return ErrServerStoped
	case err := <-errChan:
		return err
	}

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
