package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zYoma/go-url-shortener/internal/app"
	"github.com/zYoma/go-url-shortener/internal/config"
)

func main() {
	// получаем конфигурацию
	flagRunAddr, flagBaseShortURL := config.ParseFlags()

	// инициализация приложения
	application := app.New(flagRunAddr, flagBaseShortURL)

	// запустить сервис
	log.Printf("start application on %s", flagRunAddr)
	go application.Server.MustRun()

	// будем ждать сигнала остановки приложения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	//горутины выполняются пока в канал не прилетит один из ожидаемых сигналов
	sign := <-stop
	log.Println("stopping application", sign)
}
