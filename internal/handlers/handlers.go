package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

var urlStore map[string]string

func init() {
	urlStore = make(map[string]string)
}

func GenerateShortURL() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	shortURL := make([]byte, 6)
	for i := range shortURL {
		shortURL[i] = letters[rand.Intn(len(letters))]
	}
	return string(shortURL)
}

func CreateURL(w http.ResponseWriter, req *http.Request) {

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
		urlStore[shortURL] = originalURL

		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, "http://localhost:8080/%s", shortURL)
	}

}

func GetURL(w http.ResponseWriter, req *http.Request, shortURL string) {
	if shortURL == "" {
		http.Error(w, "Bad url", http.StatusBadRequest)
		return
	}

	originalURL, ok := urlStore[shortURL]
	if !ok {
		fmt.Println()
		http.NotFound(w, req)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)

}
