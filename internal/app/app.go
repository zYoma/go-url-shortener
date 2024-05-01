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

// App представляет основную структуру приложения, инкапсулирующую сервер и
// его зависимости для удобного управления.
type App struct {
	// Server - это HTTP-сервер приложения, который обрабатывает входящие запросы.
	Server     *server.HTTPServer
	stopChan   chan int64
	GRPCServer *server.GRPCServer
}

// ErrServerStoped описывает ошибку, возникающую при остановке сервера.
var ErrServerStoped = errors.New("server stoped")

// New инициализирует и возвращает новый экземпляр App, готовый к запуску.
// Эта функция принимает конфигурацию приложения и на основе её параметров
// создаёт соответствующий провайдер хранилища и сервер.
//
// cfg: параметры конфигурации приложения.
//
// Возвращает инициализированный экземпляр App и ошибку, если таковая возникла.
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
	// канал для уведомления фоновой горутины об остановке приложения
	stopChan := make(chan int64, 1)

	// создаем сервер
	httpServer := server.New(provider, cfg, stopChan)
	grpcServer := server.NewGRPC(cfg, provider)

	return &App{Server: httpServer, stopChan: stopChan, GRPCServer: grpcServer}, nil
}

// Run запускает приложение, включая HTTP-сервер и обработку сигналов
// операционной системы для корректной остановки сервера. Этот метод
// блокирует выполнение до получения сигнала остановки.
//
// Возвращает ошибку в случае неудачи запуска сервера или его остановки.
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

	// Запускаем gRPC сервер в отдельной горутине
	go func() {
		if err := s.GRPCServer.Run(); err != nil {
			errChan <- err
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-sigChan:
		// уведомляем горутину что надо остановиться
		s.stopChan <- 1
		// При получении сигнала завершения останавливаем сервер
		if err := s.Server.Shutdown(ctx); err != nil {
			return err
		}
		s.GRPCServer.Stop()
		return ErrServerStoped
	case err := <-errChan:
		return err
	}

}

// StorageConstructor в зависимости от конфигурации выбирает и возвращает
// соответствующий провайдер хранилища данных для приложения.
//
// cfg: параметры конфигурации, влияющие на выбор провайдера хранилища.
//
// Возвращает экземпляр провайдера хранилища и ошибку, если таковая возникла.
func StorageConstructor(cfg *config.Config) (storage.StorageProvider, error) {
	if cfg.DSN != "" {
		logger.Log.Sugar().Infof("провайдер - postgres")
		return postgres.New(cfg)
	}
	logger.Log.Sugar().Infof("провайдер - mem")
	return mem.New(cfg)
}
