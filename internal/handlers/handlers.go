package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/zYoma/go-url-shortener/internal/config"
	"github.com/zYoma/go-url-shortener/internal/logger"
	"github.com/zYoma/go-url-shortener/internal/models"
	"github.com/zYoma/go-url-shortener/internal/services/generator"
	"go.uber.org/zap"
)

type URLProvider interface {
	SaveURL(fullURL string, shortURL string)
	GetURL(shortURL string) (string, error)
}

// не уверен в нейминге
type HandlerService struct {
	provider URLProvider
	cfg      *config.Config
}

func New(provider URLProvider, cfg *config.Config) *HandlerService {
	return &HandlerService{provider: provider, cfg: cfg}
}

func (h *HandlerService) CreateURL(w http.ResponseWriter, req *http.Request) {
	// получаем тело запроса
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// проверяем, что тело не пустое
	originalURL := string(body)
	if originalURL == "" {
		http.Error(w, "URL cannot be empty", http.StatusBadRequest)
		return
	}

	// создаем короткую ссылку
	shortURL := generator.GenerateShortURL()

	// сохраняем ссылку в хранилище
	h.provider.SaveURL(originalURL, shortURL)

	// устанавливаем статус ответа
	w.WriteHeader(http.StatusCreated)

	// пишем ответ
	fmt.Fprintf(w, "%s/%s", h.cfg.BaseShortURL, shortURL)
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

func (h *HandlerService) CreateShortURL(w http.ResponseWriter, r *http.Request) {

	var req models.CreateShortURLRequest

	// декодируем тело запроса
	err := render.DecodeJSON(r.Body, &req)

	// если тело пустое
	if errors.Is(err, io.EOF) {
		logger.Log.Error("request body is empty")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, models.Error("empty request"))
		return
	}

	// если не удалось декодировать
	if err != nil {
		logger.Log.Error("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, models.Error("failed to decode request"))
		return
	}

	// валидируем поля
	if err := validator.New().Struct(req); err != nil {
		validateErr := err.(validator.ValidationErrors)
		logger.Log.Error("request validate error", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, models.ValidationError(validateErr))
		return
	}

	// создаем короткую ссылку
	shortURL := generator.GenerateShortURL()

	// сохраняем ссылку в хранилище
	h.provider.SaveURL(req.Url, shortURL)

	// устанавливаем статус ответа
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	// сериализуем ответ сервера
	response := models.CreateShortURLResponse{
		Result: fmt.Sprintf("%s/%s", h.cfg.BaseShortURL, shortURL),
	}

	render.JSON(w, r, response)
}

func (h *HandlerService) GetRouter() chi.Router {
	// создаем роутер
	r := chi.NewRouter()

	r.Use(handlerLogger)

	// добавляем маршруты
	r.Route("/", func(r chi.Router) {
		r.Post("/", h.CreateURL)
		r.Post("/api/shorten", h.CreateShortURL)
		r.Get("/{id}", h.GetURL)
	})

	return r
}
