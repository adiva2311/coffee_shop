package routes

import (
	"coffee_shop/config"
	"coffee_shop/controllers"
	"coffee_shop/dto"
	"coffee_shop/middlewares"
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
	g.POST("/logout", UserController.Logout, middlewares.JWTMiddleware)
	g.PATCH("/user/update", UserController.UpdateUser, middlewares.JWTMiddleware)
	g.DELETE("/user/delete", UserController.DeleteUser, middlewares.JWTMiddleware)
	g.GET("/user/detail", UserController.GetUserByID, middlewares.JWTMiddleware)

	// REFRESH TOKEN
	g.POST("/refresh-token", UserController.RefreshToken)

	// CATEGORY ROUTES
	CategoryController := controllers.NewCategoryController(db)
	g.GET("/categories", CategoryController.GetAllCategories)
	g.POST("/categories", CategoryController.CreateCategory, middlewares.JWTMiddleware)
	g.GET("/categories/:id", CategoryController.GetCategoryByID, middlewares.JWTMiddleware)
	g.PUT("/categories/:id", CategoryController.UpdateCategory, middlewares.JWTMiddleware)
	g.DELETE("/categories/:id", CategoryController.DeleteCategory, middlewares.JWTMiddleware)

	// MENU ROUTES
	MenuController := controllers.NewMenuController(db)
	g.GET("/menu", MenuController.GetAllMenus)
	g.POST("/menu", MenuController.CreateMenu, middlewares.JWTMiddleware)
	g.GET("/menu/:id", MenuController.GetMenuByID)
	g.PATCH("/menu/:id", MenuController.UpdateMenu, middlewares.JWTMiddleware)
	g.DELETE("/menu/:id", MenuController.DeleteMenu, middlewares.JWTMiddleware)
}
