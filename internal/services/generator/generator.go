package generator

import "math/rand"

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateShortURL генерирует уникальный короткий идентификатор URL.
// Используется для создания короткой ссылки из произвольного набора символов.
// Генерируемый идентификатор имеет фиксированную длину 6 символов и формируется
// из безопасного набора символов, подходящих для URL.
//
// Возвращает строку, представляющую собой короткий URL.
func GenerateShortURL() string {
	shortURL := make([]byte, 6)
	for i := range shortURL {
		shortURL[i] = letters[rand.Intn(len(letters))]
	}
	return string(shortURL)
}
