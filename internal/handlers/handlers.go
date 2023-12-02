package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zYoma/go-url-shortener/internal/services/generator"
)

type URLProvider interface {
	SaveURL(fullURL string, shortURL string)
	GetURL(shortURL string) (string, error)
}

func CreateURL(w http.ResponseWriter, req *http.Request, provider URLProvider, baseShortURL string) {
	// получаем тело запроса
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
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
	provider.SaveURL(originalURL, shortURL)

	// устанавливаем статус ответа
	w.WriteHeader(http.StatusCreated)

	// пишем ответ
	fmt.Fprintf(w, "%s/%s", baseShortURL, shortURL)
}

func GetURL(w http.ResponseWriter, req *http.Request, provider URLProvider) {
	// получаем идентификатор из пути
	shortURL := chi.URLParam(req, "id")

	// проверяем в хранилище, есть ли урл для полученного id
	originalURL, err := provider.GetURL(shortURL)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	// устанавливаем заголовок и пишем ответ
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)

}
