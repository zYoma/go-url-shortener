package main

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

func generateShortURL() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	shortURL := make([]byte, 6)
	for i := range shortURL {
		shortURL[i] = letters[rand.Intn(len(letters))]
	}
	return string(shortURL)
}

func createURL(res http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[len("/"):]

	if req.Method == http.MethodPost {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			res.Write([]byte(err.Error()))
			return
		}

		originalURL := string(body)
		if originalURL == "" {
			http.Error(res, "URL cannot be empty", http.StatusBadRequest)
			return
		}

		shortURL := generateShortURL()
		urlStore[shortURL] = originalURL

		res.WriteHeader(http.StatusCreated)

		fmt.Fprintf(res, "http://localhost:8080/%s", shortURL)
	}

	if req.Method == http.MethodGet && id != "" {
		originalURL, ok := urlStore[id]
		if !ok {
			http.NotFound(res, req)
			return
		}

		res.Header().Set("Location", originalURL)
		res.WriteHeader(http.StatusTemporaryRedirect)

	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, createURL)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
