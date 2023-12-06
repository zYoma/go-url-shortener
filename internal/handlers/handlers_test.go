package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zYoma/go-url-shortener/internal/config"
	"github.com/zYoma/go-url-shortener/internal/storage/mem"
)

func GetMockConfig() *config.Config {
	return &config.Config{
		RunAddr:      ":8080",
		BaseShortURL: "http://localhost:8080",
	}
}

func TestCreateURL(t *testing.T) {
	cfg := GetMockConfig()
	provider := mem.New()
	service := New(provider, cfg)
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
		{name: "успешный кейс", method: http.MethodPost, body: "http://ya.ru", expectedCode: http.StatusCreated, expectedBody: "http://localhost:8080/"},
		{name: "пустое тело запроса", method: http.MethodPost, body: "", expectedCode: http.StatusBadRequest, expectedBody: "URL cannot be empty"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = srv.URL
			req.SetBody(tc.body)

			resp, err := req.Send()

			require.NoError(t, err)
			assert.Equal(t, tc.expectedCode, resp.StatusCode())
			assert.Contains(t, string(resp.Body()), tc.expectedBody)
		})
	}
}

func TestGetURL(t *testing.T) {
	mockID := "sdReka"
	provider := mem.New()
	provider.SaveURL("http://ya.ru", mockID)

	cfg := GetMockConfig()
	service := New(provider, cfg)
	r := service.GetRouter()
	srv := httptest.NewServer(r)
	defer srv.Close()

	testCases := []struct {
		name         string
		method       string
		url          string
		expectedCode int
		expectedBody string
	}{
		{name: "ссылка найдена", method: http.MethodGet, url: mockID, expectedCode: http.StatusOK, expectedBody: ""},
		{name: "ссылка не найдена", method: http.MethodGet, url: "DeYqxc", expectedCode: http.StatusNotFound, expectedBody: "404 page not found"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			baseURL := srv.URL
			parsedURL, _ := url.Parse(baseURL)
			newURL := parsedURL.ResolveReference(&url.URL{Path: tc.url})
			req.URL = newURL.String()

			resp, err := req.Send()

			require.NoError(t, err)
			assert.Equal(t, tc.expectedCode, resp.StatusCode())
			if tc.expectedBody != "" {
				assert.Contains(t, string(resp.Body()), tc.expectedBody)
			}

		})
	}
}
