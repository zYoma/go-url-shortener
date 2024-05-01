package grpc_handlers

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type UserIDKeyType string

// UserIDKey - ключ для сохранения идентификатора пользователя в контексте запроса
const UserIDKey UserIDKeyType = "userID"

// AuthMiddleware - промежуточное ПО для аутентификации
func AuthMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	// Получаем user_id из заголовков метаданных
	userIDs := md.Get("user_id")
	if len(userIDs) == 0 {
		return nil, errors.New("missing user_id")
	}

	// user_id берется первое значение из слайса
	userID := userIDs[0]

	// Передаем идентификатор пользователя в контекст запроса
	ctx = context.WithValue(ctx, UserIDKey, userID)

	// Продолжаем выполнение запроса
	return handler(ctx, req)
}
