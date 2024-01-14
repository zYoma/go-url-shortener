package server

import (
	"net/http"

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
	if err != nil {
		return err
	}

	return nil
}
