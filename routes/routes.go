package routes

import (
	"coffee_shop/config"
	"coffee_shop/utils"
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func redisPing(c echo.Context) error {
	redisClient, err := config.RedisClient()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	ctx := context.Background()
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	ApiResponse := utils.ApiResponse{
		Status:  http.StatusOK,
		Message: "Coffee Shop API is running || Redis => " + pong,
	}

	return c.JSON(http.StatusOK, ApiResponse)
}

func ApiRoutes(e *echo.Echo) {
	// db, err := config.InitDB()
	// if err != nil {
	// 	log.Fatal("Failed Connect to Database")
	// }

	g := e.Group("/api/v1")

	g.GET("/health", redisPing)
}
