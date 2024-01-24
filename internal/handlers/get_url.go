package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zYoma/go-url-shortener/internal/storage/postgres"
)

func (h *HandlerService) GetURL(w http.ResponseWriter, req *http.Request) {
	// получаем идентификатор из пути
	shortURL := chi.URLParam(req, "id")

	ctx := req.Context()

	// проверяем в хранилище, есть ли урл для полученного id
	originalURL, err := h.provider.GetURL(ctx, shortURL)
	if err != nil {
		if errors.Is(err, postgres.ErrURLDeleted) {
			w.WriteHeader(http.StatusGone)
			return
		}
		http.NotFound(w, req)
		return
	}

	http.Redirect(w, req, originalURL, http.StatusTemporaryRedirect)
}
