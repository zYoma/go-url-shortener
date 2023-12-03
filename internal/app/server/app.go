package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type App struct {
	server *http.Server
}

func New(
	address string,
	router chi.Router,
) *App {
	server := &http.Server{
		Addr:    address,
		Handler: router,
	}
	return &App{
		server: server,
	}
}

func (a *App) Run() error {
	err := a.server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
