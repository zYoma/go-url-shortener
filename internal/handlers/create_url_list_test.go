package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zYoma/go-url-shortener/internal/mocks"
	"github.com/zYoma/go-url-shortener/internal/models"
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
	// Инициализация логгера, конфигурации, моковых сервисов и т.д.
	cfg := GetMockConfig()

	providerMock := new(mocks.URLProvider)
	providerMock.On("BulkSaveURL", mock.AnythingOfType("*context.valueCtx"), mock.Anything, mock.Anything).Return(nil)

	handlerService := New(providerMock, cfg)

	testURLs := []models.OriginalURL{
		{CorrelationID: "1", OriginalURL: "http://example1.com"},
		{CorrelationID: "2", OriginalURL: "http://example2.com"},
		// Добавьте больше URL для тестирования, если нужно
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
