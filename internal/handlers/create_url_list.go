package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"github.com/zYoma/go-url-shortener/internal/logger"
	"github.com/zYoma/go-url-shortener/internal/models"
	"github.com/zYoma/go-url-shortener/internal/services/generator"
	"go.uber.org/zap"
)

// CreateShortListURL обрабатывает HTTP POST запросы для создания коротких URL из списка длинных URL.
// В теле запроса ожидается JSON массив с объектами, содержащими исходные длинные URL и идентификаторы для корреляции.
// Этот метод использует стороннюю библиотеку jsoniter для эффективной работы с JSON,
// что позволяет уменьшить потребление памяти при сериализации и десериализации JSON данных.
//
// После чтения и десериализации запроса каждый URL валидируется.
// Для каждого валидного URL генерируется короткий URL, который сохраняется в хранилище с использованием
// предоставленного провайдера хранилища. В ответ клиенту отправляется JSON массив с короткими URL и их корреляционными идентификаторами.
//
// В случае неудачи при чтении тела запроса, десериализации JSON, валидации URL или сохранении в хранилище,
// клиенту отправляется соответствующий HTTP статус ошибки и описание ошибки в формате JSON.
//
// Параметры:
//
//	w http.ResponseWriter: интерфейс для отправки HTTP ответов.
//	r *http.Request: структура, представляющая HTTP запрос.
func (h *HandlerService) CreateShortListURL(w http.ResponseWriter, r *http.Request) {

	var req []models.OriginalURL

	w.Header().Set("Content-Type", "application/json")
	// в результате профилирования, выявил, что большая часть памяти тратится при сериализации,
	// использовал стороннюю библиотеку jsoniter, что позволило  уменьшить потребление памяти
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("cannot read body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, models.Error("cannot read body"))
		return
	}
	defer r.Body.Close()
	err = jsoniter.Unmarshal(body, &req)

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

	// Создаём экземпляр валидатора один раз
	validate := validator.New()

	for _, url := range req {
		// Используем уже созданный экземпляр валидатора для проверки
		if err := validate.Struct(url); err != nil {
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
	userID, err := getUserFromRequest(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = h.provider.BulkSaveURL(ctx, insertData, userID)
	if err != nil {
		render.JSON(w, r, models.Error("failed save link to db"))
		return
	}

	w.WriteHeader(http.StatusCreated)

	jsonData, err := jsoniter.Marshal(responseData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)

}
