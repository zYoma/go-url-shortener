package server

import (
	"context"
	"net/http"
	"sync"

	"github.com/zYoma/go-url-shortener/internal/config"
	"github.com/zYoma/go-url-shortener/internal/handlers"
	"github.com/zYoma/go-url-shortener/internal/storage"
)

// HTTPServer представляет собой обёртку над http.Server,
// позволяя более удобно управлять процессом его запуска и остановки.
// Включает в себя sync.WaitGroup для отслеживания активных подключений
// и корректного завершения работы сервера.
type HTTPServer struct {
	// server указывает на встроенный HTTP-сервер.
	server *http.Server

	// wg используется для отслеживания и управления параллельной обработкой
	// входящих HTTP-запросов. WaitGroup обеспечивает ожидание завершения
	// всех активных запросов перед остановкой сервера.
	wg *sync.WaitGroup
}

// New создает и возвращает новый экземпляр HTTPServer, готовый к запуску.
// Эта функция принимает провайдер URL storage.URLProvider и конфигурацию
// приложения cfg для инициализации внутренних компонентов, таких как обработчики
// и HTTP-сервер. Также инициализируется WaitGroup для контроля за горутинами,
// например, за фоновым удалением сообщений.
//
// provider: компонент для взаимодействия с хранилищем URL.
// cfg: конфигурационные параметры приложения, включая адрес запуска сервера.
//
// Возвращает указатель на инициализированный HTTPServer.
func New(
	provider storage.URLProvider,
	cfg *config.Config,
) *HTTPServer {

	// создаем сервис обработчик
	service := handlers.New(provider, cfg)

	// запускаем горутину для удаления сообщений
	var wg sync.WaitGroup
	wg.Add(1) // если нужно будет запустить несколько горутин, инкриментировать в цикле
	go service.DeleteMessages(&wg)

	// получаем роутер
	router := service.GetRouter()

	server := &http.Server{
		Addr:    cfg.RunAddr,
		Handler: router,
	}
	return &HTTPServer{
		server: server,
		wg:     &wg,
	}
}

// Run запускает HTTP-сервер на предварительно заданном адресе.
// Этот метод блокирует выполнение до тех пор, пока сервер не будет остановлен
// через вызов Shutdown или до возникновения ошибки.
//
// Возвращает ошибку, если сервер не смог запуститься или был некорректно остановлен,
// кроме ошибки http.ErrServerClosed, которая считается ожидаемым результатом
// корректной остановки сервера.
func (a *HTTPServer) Run() error {
	err := a.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Shutdown корректно останавливает HTTP-сервер, дожидаясь завершения
// всех обрабатываемых запросов и фоновых горутин, ассоциированных с сервером.
// Этот метод должен быть вызван при получении сигнала о завершении работы приложения,
// чтобы гарантировать безопасное закрытие всех ресурсов и соединений.
//
// ctx: контекст выполнения, позволяющий контролировать таймауты
// и отмену операции остановки сервера.
//
// Возвращает ошибку, если произошла ошибка при остановке сервера.
func (a *HTTPServer) Shutdown(ctx context.Context) error {
	// ждем пока все горутины завершатся
	a.wg.Wait()
	return a.server.Shutdown(ctx)
}
