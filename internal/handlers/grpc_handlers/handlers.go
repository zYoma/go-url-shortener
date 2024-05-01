package grpchandlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/zYoma/go-url-shortener/internal/config"
	"github.com/zYoma/go-url-shortener/internal/models"
	"github.com/zYoma/go-url-shortener/internal/services/generator"
	"github.com/zYoma/go-url-shortener/internal/storage"
	"github.com/zYoma/go-url-shortener/internal/storage/postgres"
	pb "github.com/zYoma/go-url-shortener/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HandlerService struct {
	provider storage.URLProvider // Интерфейс взаимодействия с хранилищем URL.
	cfg      *config.Config      // Конфигурация приложения.
	pb.UnimplementedShortenerServer
}

// New инициализирует и возвращает новый экземпляр HandlerService.
// Этот метод принимает провайдер для взаимодействия с хранилищем данных и конфигурацию приложения.
//
// provider: провайдер для взаимодействия с хранилищем данных.
// cfg: конфигурация приложения.
//
// Возвращает указатель на созданный экземпляр HandlerService.
func New(provider storage.URLProvider, cfg *config.Config) *HandlerService {
	return &HandlerService{provider: provider, cfg: cfg}
}

func (h *HandlerService) CreateShortURL(ctx context.Context, req *pb.CreateShortURLRequest) (*pb.CreateShortURLResponse, error) {
	request := models.CreateShortURLRequest{
		URL: req.GetUrl(),
	}

	if err := validator.New().Struct(request); err != nil {
		validateErr := err.(validator.ValidationErrors)
		return nil, status.Errorf(codes.InvalidArgument, "request validation error: %v", validateErr)
	}

	shortURL := generator.GenerateShortURL()

	// получаем userID из контекста
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return nil, errors.New("user ID not found in context")
	}

	if err := h.provider.SaveURL(ctx, request.URL, shortURL, userID); err != nil {
		if errors.Is(err, postgres.ErrConflict) {
			resultShortURL, _ := h.provider.GetShortURL(ctx, request.URL)
			return &pb.CreateShortURLResponse{
				Result: fmt.Sprintf("%s/%s", h.cfg.BaseShortURL, resultShortURL),
			}, status.Error(codes.AlreadyExists, "link already exists")
		}
		return nil, status.Error(codes.Internal, "failed to save link to db")
	}

	return &pb.CreateShortURLResponse{
		Result: fmt.Sprintf("%s/%s", h.cfg.BaseShortURL, shortURL),
	}, nil
}

func (h *HandlerService) GetUserURLs(ctx context.Context, req *pb.GetUserURLsRequest) (*pb.GetUserURLsResponse, error) {
	// Получаем userID из контекста
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return nil, errors.New("user ID not found in context")
	}

	userURLs, err := h.provider.GetUserURLs(ctx, h.cfg.BaseShortURL, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get links from db")
	}

	var pbUserURLs []*pb.URLs
	for _, u := range userURLs {
		url := &pb.URLs{
			ShortUrl:    u.ShortURL,
			OriginalUrl: u.OriginalURL,
		}
		pbUserURLs = append(pbUserURLs, url)
	}

	response := &pb.GetUserURLsResponse{
		Urls: pbUserURLs,
	}

	return response, nil
}

func (h *HandlerService) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	err := h.provider.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.PingResponse{Message: "OK"}, nil
}
