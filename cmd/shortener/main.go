package main

import (
	"os"
	"os/signal"
	"sync"
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
	application, err := app.New(cfg)
	if err != nil {
		panic(err)
	}

	// запускаем приложение
	if err := application.Run(); err != nil {
		panic(err)
	}

	// будем ждать сигнала остановки приложения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	//горутины выполняются пока в канал не прилетит один из ожидаемых сигналов

	// Определяем WaitGroup для ожидания завершения горутин
	var wg sync.WaitGroup
	wg.Add(1) // Добавляем горутину для ожидания

	// Горутина для ожидания сигнала SIGTERM/SIGINT
	go func() {
		defer wg.Done() // Уменьшаем счетчик горутин после завершения
		sign := <-stop
		application.Provider.Stop(cfg)
		logger.Log.Sugar().Infoln("stopping application", sign)
	}()

	// Ожидаем завершение горутины перед выходом
	wg.Wait()

	// sign := <-stop
	// application.Provider.Stop(cfg)
	// logger.Log.Sugar().Infoln("stopping application", sign)
}
