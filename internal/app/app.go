package app

import (
	"log"

	"github.com/zYoma/go-url-shortener/internal/app/server"
	"github.com/zYoma/go-url-shortener/internal/config"
	"github.com/zYoma/go-url-shortener/internal/storage/mem"
)

type App struct {
	Server *server.HTTPServer
}

func New(cfg *config.Config) *App {
	// создаем провайдер для storage
	provider := mem.New()

	// создаем сервер
	server := server.New(provider, cfg)

	return &App{Server: server}
}

func (s *App) Run() error {
	// Создание канала для ошибок
	errChan := make(chan error)

	// запустить сервис
	log.Printf("start application")

	go func() {
		errChan <- s.Server.Run()
	}()

	// Ожидание и обработка ошибки из горутины
	if err := <-errChan; err != nil {
		return err
	}

	return nil

}
