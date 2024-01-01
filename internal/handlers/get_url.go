package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

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
