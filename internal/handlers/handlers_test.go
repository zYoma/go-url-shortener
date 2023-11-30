package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateURL(t *testing.T) {
	type want struct {
		body     string
		code     int
		response string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				body:     "http://ya.ru",
				code:     201,
				response: `http://localhost:8080/`,
			},
		},
		{
			name: "empty body",
			want: want{
				body:     "",
				code:     400,
				response: `URL cannot be empty`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyReader := bytes.NewReader([]byte(tt.want.body))
			request := httptest.NewRequest(http.MethodPost, "/", bodyReader)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			CreateURL(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Contains(t, string(resBody), tt.want.response, "ответ не содержит http://localhost:8080/")
		})
	}
}

func TestGetURL(t *testing.T) {
	type want struct {
		code     int
		response string
		query    string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "url not found",
			want: want{
				query:    "sdReka",
				code:     404,
				response: `404 page not found`,
			},
		},
		{
			name: "empty query",
			want: want{
				query:    "",
				code:     400,
				response: `Bad url`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//  создадим урл чтобы было что проверять
			// bodyReader := bytes.NewReader([]byte("http://ya.ru"))
			// createRequest := httptest.NewRequest(http.MethodPost, "/", bodyReader)
			// w := httptest.NewRecorder()
			// handlers.CreateURL(w, createRequest)
			// createRes := w.Result()
			// defer createRes.Body.Close()
			// shortURL, err := io.ReadAll(createRes.Body)
			// id := strings.TrimPrefix(string(shortURL), "/")

			// пытаемся получить ранее созданные урл по id
			// httptest.NewRequest создает фейковый HTTP-запрос и не имеет возможности сохранять состояние переменных между двумя запросами.

			request := httptest.NewRequest(http.MethodGet, "/"+tt.want.query, nil)
			// создаём новый Recorder
			ww := httptest.NewRecorder()

			GetURL(ww, request, tt.want.query)

			res := ww.Result()
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Contains(t, string(resBody), tt.want.response)
		})
	}
}
