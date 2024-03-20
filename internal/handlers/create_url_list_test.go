package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zYoma/go-url-shortener/internal/config"
	"github.com/zYoma/go-url-shortener/internal/mocks"
	"github.com/zYoma/go-url-shortener/internal/models"
	"github.com/zYoma/go-url-shortener/internal/storage/mem"
)

var successBody string = `
[
	{
		"correlation_id":"d1eefc7b-2228-44af-b26d-e76711928975",
		"original_url":"http://limka.yandex/g7rl1hons/fmffjhswu5"
	},
	{
		"correlation_id":"23167760-a524-497e-a86d-3641dd8b8a56",
		"original_url":"http://lffgim2fvocaj.ru/rnyhrfkrlv/b3xloaov1o/irvbux9q9mp"
	}
]
`

var invalidURLBody string = `
[
	{
		"correlation_id":"d1eefc7b-2228-44af-b26d-e76711928975",
		"original_url":"limka.yandex/g7rl1hons/fmffjhswu5"
	},
	{
		"correlation_id":"23167760-a524-497e-a86d-3641dd8b8a56",
		"original_url":"ffgim2fvocaj.ru/rnyhrfkrlv/b3xloaov1o/irvbux9q9mp"
	}
]
`

var missingURLBody string = `
[
	{
		"original_url":"limka.yandex/g7rl1hons/fmffjhswu5"
	},
	{
		"original_url":"ffgim2fvocaj.ru/rnyhrfkrlv/b3xloaov1o/irvbux9q9mp"
	}
]
`

func TestCreateListURL(t *testing.T) {
	cfg := GetMockConfig()

	providerMock := new(mocks.URLProvider)
	providerMock.On("BulkSaveURL", mock.AnythingOfType("*context.valueCtx"), mock.Anything, mock.Anything).Return(nil)

	service := New(providerMock, cfg)
	r := service.GetRouter()
	srv := httptest.NewServer(r)
	defer srv.Close()

	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{name: "успешный кейс", method: http.MethodPost, body: successBody, expectedCode: http.StatusCreated, expectedBody: "http://localhost:8080/"},
		{name: "пустое тело запроса", method: http.MethodPost, body: "[]", expectedCode: http.StatusBadRequest, expectedBody: "empty request"},
		{name: "не верный формат url", method: http.MethodPost, body: invalidURLBody, expectedCode: http.StatusBadRequest, expectedBody: "field OriginalURL is not a valid URL"},
		{name: "не все обязательные поля переданны", method: http.MethodPost, body: missingURLBody, expectedCode: http.StatusBadRequest, expectedBody: "ield CorrelationID is a required field"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := resty.New().R()
			req.Header.Set("Accept-Encoding", "")
			req.Method = tc.method
			req.URL = fmt.Sprintf("%s/api/shorten/batch", srv.URL)
			req.SetBody(tc.body)

			resp, err := req.Send()

			require.NoError(t, err)
			assert.Equal(t, tc.expectedCode, resp.StatusCode())
			assert.Contains(t, string(resp.Body()), tc.expectedBody)

		})
	}
}

func BenchmarkCreateShortListURL(b *testing.B) {
	cfg := GetMockConfig()

	providerMock := new(mocks.URLProvider)
	providerMock.On("BulkSaveURL", mock.AnythingOfType("*context.valueCtx"), mock.Anything, mock.Anything).Return(nil)

	handlerService := New(providerMock, cfg)

	testURLs := []models.OriginalURL{
		{CorrelationID: "1", OriginalURL: "http://example1.com"},
		{CorrelationID: "2", OriginalURL: "http://example2.com"},
	}
	body, _ := json.Marshal(testURLs)

	req, err := http.NewRequest("POST", "/create-short-list-url", bytes.NewBuffer(body))
	if err != nil {
		b.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlerService.CreateShortListURL)
		handler.ServeHTTP(rr, req.WithContext(context.Background()))
	}
}

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
