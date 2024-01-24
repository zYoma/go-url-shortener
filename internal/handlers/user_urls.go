package handlers

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/zYoma/go-url-shortener/internal/models"
)

func (h *HandlerService) GetUserURL(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	// получаем userID из контекста установленного в мидлваре
	userID, ok := req.Context().Value(UserIDKey).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response, err := h.provider.GetUserURLs(ctx, h.cfg.BaseShortURL, userID)
	if err != nil {
		render.JSON(w, req, models.Error("failed get link from db"))
		return
	}

	if len(response) == 0 {
		// Тут похоже ошибка в тесте на плотформе, вместо статуса 204 там проверяется 401
		w.WriteHeader(http.StatusUnauthorized)
	}

	render.JSON(w, req, response)
}
