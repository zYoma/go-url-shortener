package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/zYoma/go-url-shortener/internal/logger"
	"github.com/zYoma/go-url-shortener/internal/models"
	"github.com/zYoma/go-url-shortener/internal/services/generator"
	"go.uber.org/zap"
)

func (h *HandlerService) CreateShortListURL(w http.ResponseWriter, r *http.Request) {

	var req []models.OriginalURL

	w.Header().Set("Content-Type", "application/json")
	err := render.DecodeJSON(r.Body, &req)

	if errors.Is(err, io.EOF) || len(req) == 0 {
		logger.Log.Error("request body is empty")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, models.Error("empty request"))
		return
	}

	if err != nil {
		logger.Log.Error("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, models.Error("failed to decode request"))
		return
	}

	for _, url := range req {
		if err := validator.New().Struct(url); err != nil {
			validateErr := err.(validator.ValidationErrors)
			logger.Log.Error("request validate error", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, models.ValidationError(validateErr))
			return
		}
	}

	var insertData []models.InsertData
	var responseData []models.ShortURL

	for _, url := range req {
		shortURL := generator.GenerateShortURL()
		insertData = append(insertData, models.InsertData{OriginalURL: url.OriginalURL, ShortURL: shortURL})
		responseData = append(responseData, models.ShortURL{CorrelationID: url.CorrelationID, ShortURL: fmt.Sprintf("%s/%s", h.cfg.BaseShortURL, shortURL)})
	}

	ctx := r.Context()

	err = h.provider.BulkSaveURL(ctx, &insertData)
	if err != nil {
		render.JSON(w, r, models.Error("failed save link to db"))
		return
	}

	w.WriteHeader(http.StatusCreated)

	render.JSON(w, r, &responseData)

}
