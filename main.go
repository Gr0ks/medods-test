package main

import (
	"medods-test/pkg/server"
	"log"
)

func main() {
	app := server.NewApp()

	if err := app.Run("8090"); err != nil {
		log.Fatalf("%s", err.Error())
	}
}





