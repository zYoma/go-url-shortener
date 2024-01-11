package models

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusError = "Error"
)

func Error(msg string) ErrorResponse {
	return ErrorResponse{
		Status: StatusError,
		Error:  msg,
	}
}

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

type CreateShortURLRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type CreateShortURLResponse struct {
	Result string `json:"result"`
}

type ShortURL struct {
	CorrelationID string `json:"correlation_id" validate:"required,correlation_id"`
	ShortURL      string `json:"short_url" validate:"required,short_url"`
}

type CreateListShortURLResponse struct {
	Result *[]ShortURL `json:"result"`
}

type OriginalURL struct {
	CorrelationID string `json:"correlation_id" validate:"required"`
	OriginalURL   string `json:"original_url" validate:"required,url"`
}

type InsertData struct {
	OriginalURL string
	ShortURL    string
}
