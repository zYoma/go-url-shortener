package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zYoma/go-url-shortener/internal/config"
)

type URLProvider interface {
	SaveURL(fullURL string, shortURL string) error
	GetURL(shortURL string) (string, error)
	Init() error
}

// не уверен в нейминге
type HandlerService struct {
	provider URLProvider
	cfg      *config.Config
}

func New(provider URLProvider, cfg *config.Config) *HandlerService {
	return &HandlerService{provider: provider, cfg: cfg}
}

func (h *HandlerService) GetURL(w http.ResponseWriter, req *http.Request) {
	// получаем идентификатор из пути
	shortURL := chi.URLParam(req, "id")

	// проверяем в хранилище, есть ли урл для полученного id
	originalURL, err := h.provider.GetURL(shortURL)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	// устанавливаем заголовок и пишем ответ
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)

}
