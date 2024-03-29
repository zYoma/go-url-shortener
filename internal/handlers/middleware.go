package handlers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/zYoma/go-url-shortener/internal/auth/jwt"
	libs "github.com/zYoma/go-url-shortener/internal/libs/gzip"
	"github.com/zYoma/go-url-shortener/internal/logger"
	"go.uber.org/zap"
)

type contextKey string

// UserIDKey представляет собой ключ, используемый для сохранения и извлечения
// идентификатора пользователя из контекста HTTP-запроса. Это позволяет передавать
// метаданные о пользователе между различными слоями приложения во время обработки запроса.
const (
	UserIDKey contextKey = "userID"
)

// ErrGetUserFromRequest описывает ошибку, возникающую при неудачной попытке
// извлечения идентификатора пользователя из контекста запроса. Это может указывать
// на отсутствие идентификатора пользователя в контексте или на проблемы с
// аутентификацией пользователя.
var ErrGetUserFromRequest = errors.New("faild get user")

func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции
		ow := w

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cw := libs.NewCompressWriter(w)
			// меняем оригинальный http.ResponseWriter на новый
			ow = cw
			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer cw.Close()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := libs.NewCompressReader(r.Body)
			if err != nil {
				logger.Log.Info("handlerLogger",
					zap.Error(err),
				)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(ow, r)

	})
}

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

func (h *HandlerService) cookieSettingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth-token")
		secret := h.cfg.TokenSecret

		if err != nil {
			// Куки нет, генерируем новый JWT токен
			tokenString, err := jwt.BuildJWTString(secret)
			if err != nil {
				// Обработка ошибки генерации токена
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Устанавливаем куку с JWT токеном
			http.SetCookie(w, &http.Cookie{
				Name:  "auth-token",
				Value: tokenString,
				Path:  "/",
			})

			// Передаем идентификатор пользователя в контекст запроса
			userID := jwt.GetUserID(tokenString, secret)
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		// Кука есть, пытаемся получить пользователя
		userID := jwt.GetUserID(cookie.Value, secret)
		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Передаем идентификатор пользователя в контекст запроса
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func getUserFromRequest(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return "", ErrGetUserFromRequest
	}
	return userID, nil
}
