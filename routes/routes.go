package routes

import (
	"coffee_shop/config"
	"coffee_shop/controllers"
	"coffee_shop/dto"
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func redisPing(c echo.Context) error {
	ctx := context.Background()

	redisClient, err := config.RedisClient()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error Connecting to Redis | " + redisClient.Ping(ctx).Err().Error(),
		})
	}

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ApiResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to ping Redis | " + redisClient.Ping(ctx).Err().Error(),
		})
	}

	ApiResponse := dto.ApiResponse{
		Status:  http.StatusOK,
		Message: "Coffee Shop API is running || Redis = " + pong,
	}

	return c.JSON(http.StatusOK, ApiResponse)
}

func ApiRoutes(e *echo.Echo) {
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed Connect to Database")
	}

	g := e.Group("/api/v1")

	g.GET("/health", redisPing)

	// USER ROUTES
	UserController := controllers.NewUserController(db)
	g.POST("/register", UserController.Register)
	g.POST("/login", UserController.Login)
	// g.PUT("/users/:user_id", UserController.UpdateUser)
	// g.DELETE("/users/:user_id", UserController.DeleteUser)
	// g.GET("/users/:user_id", UserController.GetUserByID)
}
