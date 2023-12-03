package app

import (
	"log"

	"github.com/zYoma/go-url-shortener/internal/app/server"
	"github.com/zYoma/go-url-shortener/internal/handlers"
	"github.com/zYoma/go-url-shortener/internal/storage/mem"
)

type App struct {
	Server *server.App
}

func New(address string, baseShortURL string) *App {
	// создаем провайдер для storage
	provider := mem.New()

	// создаем сервис обработчик
	service := handlers.New(provider, baseShortURL)

	// получаем роутер
	r := service.GetRouter()

	// создаем сервер
	server := server.New(address, r)

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
