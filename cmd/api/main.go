package main

import (
	"fin-flow-api/internal/bootstrap"
	"log"
)

func main() {
	app, err := bootstrap.NewApp()
	if err != nil {
		log.Fatal("failed to initialize app: ", err)
	}

	if err := app.Server.Run(); err != nil {
		log.Fatal("server error: ", err)
	}
}
