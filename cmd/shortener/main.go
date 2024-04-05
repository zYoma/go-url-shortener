package main

import (
	"errors"
	"fmt"

	"github.com/zYoma/go-url-shortener/internal/app"
	"github.com/zYoma/go-url-shortener/internal/config"
	"github.com/zYoma/go-url-shortener/internal/logger"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

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
		if errors.Is(err, app.ErrServerStoped) {
			logger.Log.Sugar().Infoln("stopping application")
			return
		}

		panic(err)
	}
}
