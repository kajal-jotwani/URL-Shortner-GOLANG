package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofibre/fibre/v2"
	"github.com/gofibre/fibre/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fibre.App) {
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)
}

func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println(err)
	}

	app := fibre.New()

	app.Use(logger.New())

	setupRoutes(app)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
