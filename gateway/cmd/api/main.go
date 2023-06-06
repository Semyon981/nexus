package main

import (
	"log"

	"github.com/Semyon981/nexus/gateway/server"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	app := server.NewApp()

	if err := app.Run("8080"); err != nil {
		log.Fatalf("%s", err.Error())
	}

}
