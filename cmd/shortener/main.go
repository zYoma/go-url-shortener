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
	flagRunAddr, flagBaseShortURL := config.ParseFlags()
	// инициализация приложения
	application := app.New(flagRunAddr, flagBaseShortURL)

	// запустить сервис
	go application.Server.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	//горутины выполняются пока в канал не прилетит один из ожидаемых сигналов
	sign := <-stop
	log.Println("stopping application", sign)
}
