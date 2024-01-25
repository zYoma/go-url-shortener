package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/render"
	"github.com/zYoma/go-url-shortener/internal/logger"
	"github.com/zYoma/go-url-shortener/internal/models"
	"go.uber.org/zap"
)

func (h *HandlerService) DeleteShortListURL(w http.ResponseWriter, req *http.Request) {
	userID, err := getUserFromRequest(req.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var listURL []string

	err = render.DecodeJSON(req.Body, &listURL)

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
	h.delChan <- models.UserListURLForDelete{UserID: userID, URLS: listURL}
	w.WriteHeader(http.StatusAccepted)
}
