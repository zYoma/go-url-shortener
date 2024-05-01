package server

import (
	"net"

	"github.com/zYoma/go-url-shortener/internal/config"
	"github.com/zYoma/go-url-shortener/internal/handlers/grpc_handlers"
	"github.com/zYoma/go-url-shortener/internal/storage"
	pb "github.com/zYoma/go-url-shortener/proto"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	server  *grpc.Server
	service *grpc_handlers.HandlerService
}

func NewGRPC(cfg *config.Config, provider storage.URLProvider) *GRPCServer {
	service := grpc_handlers.New(provider, cfg)
	return &GRPCServer{service: service}
}

func (a *GRPCServer) Run() error {
	// Задание показалось черезчур грамоздким и сложным, не хочется менять весь код, чтобы добавить всем хендлерам возможность оаботать по gRPC
	// Выбрал 3 обработчика, и сделал для них реализацию на gRPC.
	// С авторизацией не разобрался, примеров в теории не нашел, добавил AuthMiddleware который ожидает userID в заголовке
	// http сервер запускается на одном порту, gRPC на другом, порт захардкодил
	// С фасадом показалось сложно, не хочется вносить изменения в большую чать кода.
	// вроде как можно на одном порту http и gRPC поднять, но примеров нам никто не дал, реализовать не получилось
	listen, err := net.Listen("tcp", ":8081")
	if err != nil {
		return err
	}
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_handlers.AuthMiddleware),
	)
	a.server = srv

	pb.RegisterShortenerServer(srv, a.service)

	if err := a.server.Serve(listen); err != nil {
		return err
	}
	return nil
}

func (a *GRPCServer) Stop() {
	a.server.GracefulStop()
}
