package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zYoma/go-url-shortener/internal/app/server"
	"github.com/zYoma/go-url-shortener/internal/handlers"
	"github.com/zYoma/go-url-shortener/internal/storage/mem"
)

type App struct {
	Server *server.App
}

func New(address string, baseShortURL string) *App {
	provider := mem.New()
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, req *http.Request) {
			handlers.CreateURL(w, req, provider, baseShortURL)
		})
		r.Get("/{id}", func(w http.ResponseWriter, req *http.Request) {
			handlers.GetURL(w, req, provider)
		})
	})

	server := server.New(address, r)

	return &App{Server: server}
}
