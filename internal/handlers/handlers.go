package handlers

import (
	"context"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/zYoma/go-url-shortener/internal/config"
	"github.com/zYoma/go-url-shortener/internal/logger"
	"github.com/zYoma/go-url-shortener/internal/models"
	"github.com/zYoma/go-url-shortener/internal/storage"
	"go.uber.org/zap"
)

// HandlerService инкапсулирует логику обработки HTTP-запросов,
// предоставляя методы для управления короткими URL. Структура включает провайдер
// для взаимодействия с хранилищем данных, конфигурацию приложения и канал
// для асинхронного удаления списка URL, принадлежащих пользователю.
type HandlerService struct {
	provider storage.URLProvider              // Интерфейс взаимодействия с хранилищем URL.
	cfg      *config.Config                   // Конфигурация приложения.
	delChan  chan models.UserListURLForDelete // Канал для удаления списка URL.
}

// New инициализирует и возвращает новый экземпляр HandlerService.
// Этот метод принимает провайдер для взаимодействия с хранилищем данных и конфигурацию приложения.
//
// provider: провайдер для взаимодействия с хранилищем данных.
// cfg: конфигурация приложения.
//
// Возвращает указатель на созданный экземпляр HandlerService.
func New(provider storage.URLProvider, cfg *config.Config) *HandlerService {
	return &HandlerService{provider: provider, cfg: cfg, delChan: make(chan models.UserListURLForDelete, 1024)}
}

// GetRouter создает и возвращает роутер с настроенными маршрутами и middleware.
// В этом методе определяются маршруты для создания и получения коротких URL,
// а также для удаления списка URL и получения списка URL, принадлежащих пользователю.
//
// Возвращает роутер с настроенными маршрутами.
func (h *HandlerService) GetRouter() chi.Router {
	// создаем роутер
	r := chi.NewRouter()

	r.Use(handlerLogger)
	r.Use(gzipMiddleware)
	r.Use(h.cookieSettingMiddleware)

	// добавляем маршруты
	r.Route("/", func(r chi.Router) {
		r.Post("/", h.CreateURL)
		r.Post("/api/shorten", h.CreateShortURL)
		r.Get("/{id}", h.GetURL)
		r.Get("/ping", h.Ping)
		r.Post("/api/shorten/batch", h.CreateShortListURL)
		r.Get("/api/user/urls", h.GetUserURL)
		r.Delete("/api/user/urls", h.DeleteShortListURL)
		r.Get("/api/internal/stats", h.GetStats)
	})

	return r
}

// DeleteMessages постоянно слушает канал delChan для асинхронного получения и удаления
// списков URL из хранилища. Метод использует таймер для периодического удаления
// накопившихся списков URL и прекращает работу при получении сигнала завершения.
//
// wg *sync.WaitGroup: группа ожидания для синхронизации завершения горутины.
func (h *HandlerService) DeleteMessages(wg *sync.WaitGroup, stopChan chan int64) {
	defer wg.Done()

	// будем сохранять сообщения, накопленные за последние 10 секунд
	ticker := time.NewTicker(10 * time.Second)

	var messages []models.UserListURLForDelete

	for {
		select {
		case msg := <-h.delChan:
			// пришло новое сообщение, добавляем в список на удаление
			messages = append(messages, msg)
			if len(messages) >= 100 {
				// Если в списке накопилось 100 сообщений, запускаем удаление
				h.saveMessages(&messages)
			}
		case <-ticker.C:
			// сработал таймер, запускаем удаление
			h.saveMessages(&messages)
		case <-stopChan:
			// сигнал остановки приложение
			h.saveMessages(&messages)
			return
		}
	}
}

// saveMessages удаляет списки URL, указанные в messages, из хранилища.
// Этот внутренний метод вызывается из DeleteMessages для фактического удаления данных.
// После успешного удаления список messages очищается.
//
// messages *[]models.UserListURLForDelete: указатель на список сообщений для удаления.
func (h *HandlerService) saveMessages(messages *[]models.UserListURLForDelete) {
	if len(*messages) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := h.provider.DeleteListURL(ctx, *messages)
	if err != nil {
		logger.Log.Debug("cannot save messages", zap.Error(err))
		return
	}

	// Очищаем сообщения после сохранения
	*messages = nil
}
