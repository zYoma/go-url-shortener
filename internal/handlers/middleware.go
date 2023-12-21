package handlers

import (
	"net/http"
	"time"

	"github.com/zYoma/go-url-shortener/internal/logger"
	"go.uber.org/zap"
)

func handlerLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создание ResponseRecorder для перехвата ответа
		recorder := &responseRecorder{w, 0, 0}

		next.ServeHTTP(recorder, r)

		duration := time.Since(start)

		logger.Log.Info("handlerLogger",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Duration("duration", duration),
			zap.Int("status", recorder.status),
			zap.Int64("size", recorder.size),
		)

	})
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	size   int64
}

// Переопределение WriteHeader для сохранения реального статуса ответа
func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Переопределение Write для сохранения размера ответа
func (r *responseRecorder) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.size += int64(size)
	return size, err
}
