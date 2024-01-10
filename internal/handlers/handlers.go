package handlers

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/zYoma/go-url-shortener/internal/config"
)

type URLProvider interface {
	SaveURL(ctx context.Context, fullURL string, shortURL string) error
	GetURL(ctx context.Context, shortURL string) (string, error)
	Init() error
	Ping(ctx context.Context) error
}

// не уверен в нейминге
type HandlerService struct {
	provider URLProvider
	cfg      *config.Config
}

func New(provider URLProvider, cfg *config.Config) *HandlerService {
	return &HandlerService{provider: provider, cfg: cfg}
}

func (h *HandlerService) GetRouter() chi.Router {
	// создаем роутер
	r := chi.NewRouter()

	r.Use(handlerLogger)
	r.Use(gzipMiddleware)

	// добавляем маршруты
	r.Route("/", func(r chi.Router) {
		r.Post("/", h.CreateURL)
		r.Post("/api/shorten", h.CreateShortURL)
		r.Get("/{id}", h.GetURL)
		r.Get("/ping", h.Ping)
	})

	return r
}
