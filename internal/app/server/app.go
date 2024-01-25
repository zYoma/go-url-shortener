package server

import (
	"context"
	"net/http"
	"time"

	"github.com/zYoma/go-url-shortener/internal/config"
	"github.com/zYoma/go-url-shortener/internal/handlers"
	"github.com/zYoma/go-url-shortener/internal/storage"
)

type HTTPServer struct {
	server *http.Server
}

func New(
	provider storage.URLProvider,
	cfg *config.Config,
) *HTTPServer {

	// создаем сервис обработчик
	service := handlers.New(provider, cfg)

	// запускаем горутину для удаления сообщений
	go service.DeleteMessages()

	// получаем роутер
	router := service.GetRouter()

	server := &http.Server{
		Addr:    cfg.RunAddr,
		Handler: router,
	}
	return &HTTPServer{
		server: server,
	}
}

func (a *HTTPServer) Run() error {
	err := a.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (a *HTTPServer) Shutdown(ctx context.Context) error {
	// ждем пока фоновые таски завершатся
	time.Sleep(5 * time.Second)
	return a.server.Shutdown(ctx)
}
