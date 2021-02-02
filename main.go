package main

import (
	"medods-test/pkg/server"
)

func main() {
	app := server.NewApp()

	app.Run("8090")
}





