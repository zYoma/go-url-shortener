package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/zYoma/go-url-shortener/internal/app"
	"github.com/zYoma/go-url-shortener/internal/config"
	"github.com/zYoma/go-url-shortener/internal/logger"
)

func main() {
	// получаем конфигурацию
	cfg := config.GetConfig()

	// инициализируем логер
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		panic(err)
	}

	// инициализация приложения
	application := app.New(cfg)

	// запускаем приложение
	if err := application.Run(); err != nil {
		panic(err)
	}

	// будем ждать сигнала остановки приложения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	//горутины выполняются пока в канал не прилетит один из ожидаемых сигналов
	sign := <-stop
	logger.Log.Sugar().Infoln("stopping application", sign)
}
