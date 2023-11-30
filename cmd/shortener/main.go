package main

import (
	"net/http"
	"strings"

	"github.com/zYoma/go-url-shortener/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Обработка GET запросов по адресу /{любая строка}
			path := r.URL.Path
			if path != "/" {
				shortURL := strings.TrimPrefix(path, "/")
				handlers.GetURL(w, r, shortURL)
				return
			}
		} else if r.Method == http.MethodPost {
			// Обработка POST запросов по адресу /
			handlers.CreateURL(w, r)
			return
		}

		http.NotFound(w, r)
	})

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
