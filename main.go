package main

import (
	"coffee_shop/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()

	// ROUTES
	routes.ApiRoutes(e)

	// Start server
	api_host := os.Getenv("API_HOST")
	api_port := os.Getenv("API_PORT")
	e.Logger.Fatal(e.Start(api_host + ":" + api_port))
}
