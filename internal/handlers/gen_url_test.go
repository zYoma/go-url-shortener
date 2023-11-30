package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zYoma/go-url-shortener/internal/handlers"
)

func TestGenerateShortURL(t *testing.T) {
	for i := 0; i < 5; i++ {
		url := handlers.GenerateShortURL()
		// проверяем что получаем строку
		assert.IsType(t, "", url)
		// проверяем, что строка содержит 6 символов
		assert.Len(t, url, 6)

	}
}
