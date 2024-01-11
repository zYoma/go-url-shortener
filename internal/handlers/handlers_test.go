package handlers

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zYoma/go-url-shortener/internal/config"
	"github.com/zYoma/go-url-shortener/internal/mocks"
	"github.com/zYoma/go-url-shortener/internal/storage/mem"
	"github.com/zYoma/go-url-shortener/internal/storage/postgres"
)

func GetMockConfig() *config.Config {
	return &config.Config{
		RunAddr:      ":8080",
		BaseShortURL: "http://localhost:8080",
		StorageFile:  "/tmp/short-url-db.json",
	}
}

func TestCreateURL(t *testing.T) {
	cfg := GetMockConfig()

	providerMock := new(mocks.URLProvider)
	// Настройка поведения мока для метода SaveURL
	providerMock.On("SaveURL", mock.AnythingOfType("*context.valueCtx"), mock.Anything, mock.Anything).Return(nil)

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
		{name: "успешный кейс", method: http.MethodPost, body: "http://ya.ru", expectedCode: http.StatusCreated, expectedBody: "http://localhost:8080/"},
		{name: "пустое тело запроса", method: http.MethodPost, body: "", expectedCode: http.StatusBadRequest, expectedBody: "URL cannot be empty"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := resty.New().R()
			req.Header.Set("Accept-Encoding", "")
			req.Method = tc.method
			req.URL = srv.URL
			req.SetBody(tc.body)

			resp, err := req.Send()

			require.NoError(t, err)
			assert.Equal(t, tc.expectedCode, resp.StatusCode())
			assert.Contains(t, string(resp.Body()), tc.expectedBody)

			// Проверка вызовов методов
			providerMock.AssertCalled(t, "SaveURL", mock.AnythingOfType("*context.valueCtx"), mock.Anything, mock.Anything)
		})
	}
}

func TestGetURL(t *testing.T) {
	mockID := "sdReka"
	cfg := GetMockConfig()
	providerMock := new(mocks.URLProvider)
	// задаем поведение для аргумента mockID и всех остальных
	providerMock.On("GetURL", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("string")).Return(
		func(ctx context.Context, shortURL string) string {
			url := ""
			if shortURL == mockID {
				url = "https://httpbin.org/get"
			}
			return url
		}, func(ctx context.Context, shortURL string) error {
			if shortURL != mockID {
				return mem.ErrURLNotFound
			}
			return nil
		})

	service := New(providerMock, cfg)
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
			req.Header.Set("Accept-Encoding", "")
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

func TestCreateShortURL(t *testing.T) {
	cfg := GetMockConfig()
	providerMock := new(mocks.URLProvider)
	providerMock.On("SaveURL", mock.AnythingOfType("*context.valueCtx"), mock.Anything, mock.Anything).Return(
		func(ctx context.Context, fullURL string, shortURL string) error {
			if fullURL == "http://mail.ru" {
				return postgres.ErrConflict
			}
			return nil
		},
	)
	providerMock.On("GetShortURL", mock.AnythingOfType("*context.valueCtx"), mock.Anything).Return("conflict", nil)

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
		{name: "успешный кейс", method: http.MethodPost, body: `{"url": "http://ya.ru"}`, expectedCode: http.StatusCreated, expectedBody: "http://localhost:8080/"},
		{name: "пустое тело запроса", method: http.MethodPost, body: "", expectedCode: http.StatusBadRequest, expectedBody: "empty request"},
		{name: "невалидный json", method: http.MethodPost, body: `{"url": "http://ya.ru",}`, expectedCode: http.StatusBadRequest, expectedBody: "failed to decode request"},
		{name: "невалидный url", method: http.MethodPost, body: `{"url": "ya.ru"}`, expectedCode: http.StatusBadRequest, expectedBody: "is not a valid URL"},
		{name: "не передан url", method: http.MethodPost, body: `{}`, expectedCode: http.StatusBadRequest, expectedBody: "URL is a required field"},
		{name: "url уже существует в БД", method: http.MethodPost, body: `{"url": "http://mail.ru"}`, expectedCode: http.StatusConflict, expectedBody: "conflict"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := resty.New().R()
			req.Header.Set("Accept-Encoding", "")
			req.Method = tc.method
			url := fmt.Sprintf("%s/api/shorten", srv.URL)
			req.URL = url
			req.SetBody(tc.body)

			resp, err := req.Send()

			require.NoError(t, err)
			contentType := resp.Header().Get("Content-Type")
			assert.Equal(t, "application/json", contentType)
			assert.Equal(t, tc.expectedCode, resp.StatusCode())
			if tc.expectedBody != "" {
				assert.Contains(t, string(resp.Body()), tc.expectedBody)
			}
		})
	}
}

func TestGzipCompression(t *testing.T) {
	cfg := GetMockConfig()
	providerMock := new(mocks.URLProvider)
	providerMock.On("SaveURL", mock.AnythingOfType("*context.valueCtx"), mock.Anything, mock.Anything).Return(nil)
	service := New(providerMock, cfg)
	r := service.GetRouter()
	srv := httptest.NewServer(r)
	defer srv.Close()

	requestBody := `{}`

	// ожидаемое содержимое тела ответа при успешном запросе
	successBody := `{
		"status": "Error",
		"error": "field URL is a required field"
	}`

	t.Run("sends_gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte(requestBody))

		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		url := fmt.Sprintf("%s/api/shorten", srv.URL)
		r := httptest.NewRequest("POST", url, buf)
		r.RequestURI = ""
		r.Header.Set("Content-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		rr := bytes.NewReader(b)
		gr, _ := gzip.NewReader(rr)
		defer gr.Close()
		// Чтение и распаковка данных
		var result bytes.Buffer
		io.Copy(&result, gr)

		require.JSONEq(t, successBody, result.String())
	})

	t.Run("accepts_gzip", func(t *testing.T) {
		buf := bytes.NewBufferString(requestBody)
		url := fmt.Sprintf("%s/api/shorten", srv.URL)
		r := httptest.NewRequest("POST", url, buf)
		r.RequestURI = ""
		r.Header.Set("Accept-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		defer resp.Body.Close()

		zr, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)

		b, err := io.ReadAll(zr)
		require.NoError(t, err)

		require.JSONEq(t, successBody, string(b))
	})
}
