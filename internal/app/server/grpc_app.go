package server

import (
	"net"

	"github.com/zYoma/go-url-shortener/internal/config"
	grpchandlers "github.com/zYoma/go-url-shortener/internal/handlers/grpc_handlers"
	"github.com/zYoma/go-url-shortener/internal/storage"
	pb "github.com/zYoma/go-url-shortener/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	server  *grpc.Server
	service *grpchandlers.HandlerService
}

func NewGRPC(cfg *config.Config, provider storage.URLProvider) *GRPCServer {
	service := grpchandlers.New(provider, cfg)
	return &GRPCServer{service: service}
}

func (a *GRPCServer) Run() error {
	// Задание показалось сложным и не своевременным, не хочется в самом конце менять весь код, чтобы добавить всем хендлерам возможность работать по gRPC
	// С фасадом показалось сложно, не хочется вносить изменения в большую чать кода.
	// Выбрал 3 обработчика, и сделал для них реализацию на gRPC.
	// С авторизацией не разобрался, примеров в теории не нашел, добавил AuthMiddleware который ожидает userID в заголовке
	// http сервер запускается на одном порту, gRPC на другом, порт захардкодил
	// вроде как можно на одном порту http и gRPC поднять, но примеров нам никто не дал, реализовать не получилось
	listen, err := net.Listen("tcp", ":8081")
	if err != nil {
		return err
	}
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpchandlers.AuthMiddleware),
	)
	a.server = srv

	pb.RegisterShortenerServer(srv, a.service)
	reflection.Register(srv)

	if err := a.server.Serve(listen); err != nil {
		return err
	}
	return nil
}

func (a *GRPCServer) Stop() {
	a.server.GracefulStop()
}
