package bootstrap

import (
	"os"

	httptransport "fin-flow-api/internal/interfaces/http"
)

type App struct {
	Server *httptransport.Server
}

func NewApp() (*App, error) {
	cfg := httptransport.Config{
		Addr: os.Getenv("PORT"),
	}

	srv := httptransport.NewServer(cfg)

	return &App{Server: srv}, nil
}