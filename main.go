package main

import (
	"flag"

	"github.com/agonist/hotel-reservation/api"
	"github.com/gofiber/fiber/v2"
)

func main(){
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	app := fiber.New()
	apiV1 := app.Group("/api/v1")

	apiV1.Get("/user", api.HandleListUsers)
	apiV1.Get("user/:id", api.HandleGetUser)

	app.Listen(*listenAddr)
}
