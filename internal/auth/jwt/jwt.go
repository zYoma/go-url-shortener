package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/zYoma/go-url-shortener/internal/logger"
)

// Claims определяет структуру утверждений (claims), используемых в JWT.
// Включает в себя стандартные зарегистрированные утверждения и пользовательское утверждение UserID.
type Claims struct {
	jwt.RegisteredClaims        // Встроенные стандартные утверждения JWT.
	UserID               string // Уникальный идентификатор пользователя.
}

// TokenExp задаёт время жизни токена.
const TokenExp = time.Hour * 3

// ErrCreateUUID определяет ошибку, возникающую при неудачной попытке генерации UUID.
var ErrCreateUUID = errors.New("create uuid")

// ErrCreateToken определяет ошибку, возникающую при неудачной попытке создания JWT.
var ErrCreateToken = errors.New("create token")

// generateUUID генерирует и возвращает уникальный идентификатор (UUID).
func generateUUID() (string, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return newUUID.String(), nil
}

// BuildJWTString создаёт JWT токен для идентификатора пользователя, генерируемого функцией generateUUID,
// и возвращает его в виде строки. Использует секрет для подписи токена.
//
// secret: секретный ключ, используемый для подписи токена.
//
// Возвращает строку, содержащую JWT токен, и ошибку, если таковая возникла.
func BuildJWTString(secret string) (string, error) {
	uuid, err := generateUUID()
	if err != nil {
		logger.Log.Sugar().Errorf("Не удалось сгенерировать uuid: %s", err)
		return "", ErrCreateUUID
	}

	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		// собственное утверждение
		UserID: uuid,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		logger.Log.Sugar().Errorf("Не удалось создать токен: %s", err)
		return "", ErrCreateToken
	}

	// возвращаем строку токена
	return tokenString, nil
}

// GetUserID извлекает идентификатор пользователя из JWT токена,
// используя указанный секретный ключ для проверки подписи токена.
//
// tokenString: строка, содержащая JWT токен.
// secret: секретный ключ, используемый для проверки подписи токена.
//
// Возвращает идентификатор пользователя из токена или пустую строку, если токен невалиден или произошла ошибка.
func GetUserID(tokenString string, secret string) string {

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		})
	if err != nil {
		logger.Log.Sugar().Errorf("Неожиданный метод подписи: %s", err)
		return ""
	}

	if !token.Valid {
		logger.Log.Sugar().Errorf("Токен не валидный")
		return ""
	}

	return claims.UserID
}
