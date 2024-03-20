package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zYoma/go-url-shortener/internal/storage/postgres"
)

// GetURL обрабатывает HTTP-запросы для перенаправления пользователя по короткой ссылке.
// Метод извлекает идентификатор короткой ссылки из URL-параметра запроса, выполняет поиск
// оригинального URL в хранилище по данному идентификатору и, в случае успеха, перенаправляет пользователя
// по найденному оригинальному URL.
//
// В случае, если URL был удалён, клиенту возвращается HTTP-статус 410 (Gone),
// указывающий на то, что ресурс был удалён и более недоступен.
// Если соответствующий оригинальный URL не найден, возвращается статус 404 (Not Found).
//
// Параметры:
//
//	w http.ResponseWriter: интерфейс для отправки HTTP ответов и перенаправлений.
//	req *http.Request: структура, представляющая HTTP запрос и содержащая параметры URL.
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
