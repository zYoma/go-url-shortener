package handlers

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/zYoma/go-url-shortener/internal/config"
	"github.com/zYoma/go-url-shortener/internal/logger"
	"github.com/zYoma/go-url-shortener/internal/models"
	"github.com/zYoma/go-url-shortener/internal/storage"
	"go.uber.org/zap"
)

type HandlerService struct {
	provider storage.URLProvider
	cfg      *config.Config
	delChan  chan models.UserListURLForDelete
}

func New(provider storage.URLProvider, cfg *config.Config) *HandlerService {
	return &HandlerService{provider: provider, cfg: cfg, delChan: make(chan models.UserListURLForDelete, 1024)}
}

func (h *HandlerService) GetRouter() chi.Router {
	// создаем роутер
	r := chi.NewRouter()

	r.Use(handlerLogger)
	r.Use(gzipMiddleware)
	r.Use(h.cookieSettingMiddleware)

	// добавляем маршруты
	r.Route("/", func(r chi.Router) {
		r.Post("/", h.CreateURL)
		r.Post("/api/shorten", h.CreateShortURL)
		r.Get("/{id}", h.GetURL)
		r.Get("/ping", h.Ping)
		r.Post("/api/shorten/batch", h.CreateShortListURL)
		r.Get("/api/user/urls", h.GetUserURL)
		r.Delete("/api/user/urls", h.DeleteShortListURL)
	})

	return r
}

// deleteMessages постоянно удаляет несколько сообщений из хранилища с определённым интервалом
func (h *HandlerService) DeleteMessages() {
	// будем сохранять сообщения, накопленные за последние 10 секунд
	ticker := time.NewTicker(100 * time.Second)

	var messages []models.UserListURLForDelete

	// Канал для перехвата сигналов завершения работы
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case msg := <-h.delChan:
			// пришло новое сообщение, добавляем в список на удаление
			messages = append(messages, msg)
		case <-ticker.C:
			// сработал таймер, запускаем удаление
			h.saveMessages(&messages)
		case <-sigChan:
			// сигнал остановки приложение
			h.saveMessages(&messages)
			return
		}
	}
}

func (h *HandlerService) saveMessages(messages *[]models.UserListURLForDelete) {
	if len(*messages) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := h.provider.DeleteListURL(ctx, *messages)
	if err != nil {
		logger.Log.Debug("cannot save messages", zap.Error(err))
		return
	}

	// Очищаем сообщения после сохранения
	*messages = nil
}
