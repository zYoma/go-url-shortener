package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type URLProvider interface {
	SaveUrl(fullURL string, shortURL string)
	GetUrl(shortURL string) (string, error)
}

func GenerateShortURL() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	shortURL := make([]byte, 6)
	for i := range shortURL {
		shortURL[i] = letters[rand.Intn(len(letters))]
	}
	return string(shortURL)
}

func CreateURL(w http.ResponseWriter, req *http.Request, provider URLProvider) {

	if req.Method == http.MethodPost {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		originalURL := string(body)
		if originalURL == "" {
			http.Error(w, "URL cannot be empty", http.StatusBadRequest)
			return
		}

		shortURL := GenerateShortURL()
		provider.SaveUrl(originalURL, shortURL)

		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, "http://localhost:8080/%s", shortURL)
	}

}

func GetURL(w http.ResponseWriter, req *http.Request, provider URLProvider) {
	shortURL := chi.URLParam(req, "id")

	originalURL, err := provider.GetUrl(shortURL)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)

}
