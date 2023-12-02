package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zYoma/go-url-shortener/internal/storage/mem"
)

func TestCreateURL(t *testing.T) {
	provider := mem.New()
	handler := func(w http.ResponseWriter, req *http.Request) {
		CreateURL(w, req, provider)
	}
	srv := httptest.NewServer(http.HandlerFunc(handler))
	defer srv.Close()

	testCases := []struct {
		method       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{method: http.MethodPost, body: "http://ya.ru", expectedCode: http.StatusCreated, expectedBody: "http://localhost:8080/"},
		{method: http.MethodPost, body: "", expectedCode: http.StatusBadRequest, expectedBody: "URL cannot be empty"},
	}
	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
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

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", func(w http.ResponseWriter, req *http.Request) {
			GetURL(w, req, provider)
		})
	})
	srv := httptest.NewServer(router)
	defer srv.Close()

	testCases := []struct {
		method       string
		url          string
		expectedCode int
		expectedBody string
	}{
		{method: http.MethodGet, url: mockID, expectedCode: http.StatusOK, expectedBody: ""},
		{method: http.MethodGet, url: "DeYqxc", expectedCode: http.StatusNotFound, expectedBody: "404 page not found"},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
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
