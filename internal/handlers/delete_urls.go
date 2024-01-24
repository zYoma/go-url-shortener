package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/zYoma/go-url-shortener/internal/logger"
	"github.com/zYoma/go-url-shortener/internal/models"
	"go.uber.org/zap"
)

func (h *HandlerService) DeleteShortListURL(w http.ResponseWriter, req *http.Request) {
	userID, _ := req.Context().Value(UserIDKey).(string)
	var listURL []string

	err := render.DecodeJSON(req.Body, &listURL)

	if errors.Is(err, io.EOF) {
		logger.Log.Error("request body is empty")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, req, models.Error("empty request"))
		return
	}
	if err != nil {
		logger.Log.Error("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, req, models.Error("failed to decode request"))
		return
	}

	go func() {
		// Создаем новый контекст с тайм-аутом
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Выполняем операцию с новым контекстом
		if err := h.provider.DeleteListURL(ctx, listURL, userID); err != nil {
			logger.Log.Error("Error deleting URLs", zap.Error(err))
			return
		}
	}()
	w.WriteHeader(http.StatusAccepted)
}
