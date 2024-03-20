package models

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ErrorResponse представляет собой структуру ответа на HTTP-запрос,
// содержащую информацию об ошибке. Используется для унификации формата
// сообщений об ошибках, возвращаемых клиенту.
type ErrorResponse struct {
	Status string `json:"status"`          // Статус ответа, указывающий на наличие ошибки.
	Error  string `json:"error,omitempty"` // Описание ошибки.
}

// StatusError - константа, обозначающая статус ошибки в ответе.
const (
	StatusError = "Error"
)

// Error создает объект ErrorResponse с указанным сообщением об ошибке.
func Error(msg string) ErrorResponse {
	return ErrorResponse{
		Status: StatusError,
		Error:  msg,
	}
}

// ValidationError создает объект ErrorResponse на основе ошибок валидации,
// предоставленных пакетом validator. Сообщения об ошибках агрегируются и форматируются
// для удобного отображения клиенту.
func ValidationError(errs validator.ValidationErrors) ErrorResponse {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return ErrorResponse{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}

// CreateShortURLRequest описывает структуру входящего запроса на создание короткой ссылки.
// Содержит URL, который требуется сократить.
type CreateShortURLRequest struct {
	URL string `json:"url" validate:"required,url"` // URL для сокращения, должен быть валидным и указан.
}

// CreateShortURLResponse описывает структуру ответа на запрос создания короткой ссылки.
// Содержит результат в виде сокращенного URL.
type CreateShortURLResponse struct {
	Result string `json:"result"` // Результат сокращения - короткий URL.
}

// ShortURL представляет собой структуру, содержащую информацию о коротком URL и связанном с ним идентификаторе корреляции.
type ShortURL struct {
	CorrelationID string `json:"correlation_id" validate:"required"` // Идентификатор для корреляции запроса и ответа.
	ShortURL      string `json:"short_url" validate:"required"`      // Сокращенный URL.
}

// OriginalURL описывает структуру с исходным URL и связанным идентификатором корреляции.
type OriginalURL struct {
	CorrelationID string `json:"correlation_id" validate:"required"`   // Идентификатор для корреляции.
	OriginalURL   string `json:"original_url" validate:"required,url"` // Исходный URL, который был сокращен.
}

// InsertData содержит данные для вставки в хранилище: оригинальный и короткий URL.
type InsertData struct {
	OriginalURL string // Исходный URL.
	ShortURL    string // Сокращенный URL.
}

// UserURLS описывает структуру данных, возвращаемую пользователю, содержащую короткий и исходный URL.
type UserURLS struct {
	ShortURL    string `json:"short_url"`    // Короткий URL.
	OriginalURL string `json:"original_url"` // Исходный URL.
}

// UserListURLForDelete описывает структуру данных для запроса на удаление списка URL пользователя.
type UserListURLForDelete struct {
	UserID string   // Идентификатор пользователя.
	URLS   []string // Список URL для удаления.
}
