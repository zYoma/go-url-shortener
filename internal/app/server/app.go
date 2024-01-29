package server

import (
	"context"
	"net/http"
	"sync"

	"github.com/zYoma/go-url-shortener/internal/config"
	"github.com/zYoma/go-url-shortener/internal/handlers"
	"github.com/zYoma/go-url-shortener/internal/storage"
)

type HTTPServer struct {
	server *http.Server
	wg     *sync.WaitGroup
}

func New(
	provider storage.URLProvider,
	cfg *config.Config,
) *HTTPServer {

	// создаем сервис обработчик
	service := handlers.New(provider, cfg)

	// запускаем горутину для удаления сообщений
	var wg sync.WaitGroup
	wg.Add(1) // если нужно будет запустить несколько горутин, инкриментировать в цикле
	go service.DeleteMessages(&wg)

	// получаем роутер
	router := service.GetRouter()

	server := &http.Server{
		Addr:    cfg.RunAddr,
		Handler: router,
	}
	return &HTTPServer{
		server: server,
		wg:     &wg,
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
	// ждем пока все горутины завершатся
	a.wg.Wait()
	return a.server.Shutdown(ctx)
}
