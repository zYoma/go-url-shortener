package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/zYoma/go-url-shortener/internal/logger"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

const TokenExp = time.Hour * 3

var ErrCreateUUID = errors.New("error create uuid")
var ErrCreateToken = errors.New("error create token")

func generateUUID() (string, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return newUUID.String(), nil
}

// BuildJWTString создаёт токен и возвращает его в виде строки.
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
