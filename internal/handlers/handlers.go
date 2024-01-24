package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/zYoma/go-url-shortener/internal/config"
	"github.com/zYoma/go-url-shortener/internal/storage"
)

type HandlerService struct {
	provider storage.URLProvider
	cfg      *config.Config
}

func New(provider storage.URLProvider, cfg *config.Config) *HandlerService {
	return &HandlerService{provider: provider, cfg: cfg}
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
	})

	return r
}
