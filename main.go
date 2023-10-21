package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/tepavcevic/hotel-reservation/api"
)

func main() {
	listenAddr := flag.String("listenAddr", ":8080", "Port for our server")
	flag.Parse()
	app := fiber.New()
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/user", api.HandleGetUsers)
	apiv1.Get("/user/:id", api.HandleGetUserById)
	app.Listen(*listenAddr)
}
