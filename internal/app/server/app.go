package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type App struct {
	server *http.Server
	host   string
	port   int
	router chi.Router
}

func New(
	host string,
	port int,
	router chi.Router,
) *App {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: router,
	}
	return &App{
		server: server,
		host:   host,
		port:   port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	err := a.server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
