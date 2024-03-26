package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortURL(t *testing.T) {
	for i := 0; i < 5; i++ {
		url := GenerateShortURL()
		// проверяем что получаем строку
		assert.IsType(t, "", url)
		// проверяем, что строка содержит 6 символов
		assert.Len(t, url, 6)

	}
}

func BenchmarkGenerateShortURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateShortURL()
	}
}
