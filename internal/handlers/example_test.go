package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/zYoma/go-url-shortener/internal/config"
	"github.com/zYoma/go-url-shortener/internal/models"
	"github.com/zYoma/go-url-shortener/internal/storage/mem"
)

func ExampleHandlerService_CreateShortListURL() {
	// получаем конфигурацию
	cfg := config.GetConfig()

	// создаем провайдер для storage
	provider, _ := mem.New(cfg)

	// Создание экземпляра HandlerService.
	h := New(provider, cfg)

	// Подготовка данных запроса.
	originalURLs := []models.OriginalURL{
		{OriginalURL: "http://example.com/longurl1", CorrelationID: "1"},
		{OriginalURL: "http://example.com/longurl2", CorrelationID: "2"},
	}
	body, _ := jsoniter.Marshal(originalURLs)

	// Создание HTTP запроса с данными.
	req, _ := http.NewRequest("POST", "/api/shorten/batch", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Настройка роутера с обработчиком.
	r := chi.NewRouter()
	r.Post("/api/shorten/batch", h.CreateShortListURL)

	// Вызов обработчика через роутер.
	r.ServeHTTP(w, req)

	fmt.Println(w.Code) // Вывод статуса ответа.

	// Пример вывода: 201
}

func ExampleHandlerService_DeleteShortListURL() {
	// получаем конфигурацию
	cfg := config.GetConfig()

	// создаем провайдер для storage
	provider, _ := mem.New(cfg)

	// Создание экземпляра HandlerService.
	h := New(provider, cfg)
	h.delChan = make(chan models.UserListURLForDelete, 1)

	// Подготовка данных запроса: список коротких URL для удаления.
	listURL := []string{"short1", "short2"}
	body, _ := jsoniter.Marshal(listURL)

	// Создание HTTP запроса с данными.
	req, _ := http.NewRequest("DELETE", "/api/user/urls", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Установка заголовка для имитации аутентификации пользователя (в вашем случае механизм может отличаться).
	req.Header.Set("Authorization", "Bearer user-token")

	// Настройка роутера с обработчиком.
	r := chi.NewRouter()
	r.Delete("/api/user/urls", h.DeleteShortListURL)

	// Вызов обработчика через роутер.
	r.ServeHTTP(w, req)

	fmt.Println(w.Code) // Вывод статуса ответа.

	// Пример вывода: 202
}
