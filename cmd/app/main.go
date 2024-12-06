// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"log"
	"os"

	"github.com/stwrtrio/movie-festival/config"
	"github.com/stwrtrio/movie-festival/internal/controllers"
	"github.com/stwrtrio/movie-festival/internal/middlewares"
	"github.com/stwrtrio/movie-festival/internal/repositories"
	"github.com/stwrtrio/movie-festival/internal/routes"
	"github.com/stwrtrio/movie-festival/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/stwrtrio/movie-festival/docs" // Import the generated docs package
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to database
	config.InitDB()
	defer config.DB.Close()

	// Initialize Redis
	config.InitRedis()
	defer config.RedisClient.Close()

	// Initialize Echo
	e := echo.New()
	e.Validator = &middlewares.CustomValidator{Validator: validator.New()}

	// Dependency Injection
	// Repository
	movieRepo := repositories.NewMovieRepository(config.DB)
	userRepo := repositories.NewUserRepository(config.DB)

	// Service
	movieService := services.NewMovieService(movieRepo, config.RedisClient)
	userService := services.NewUserService(userRepo, config.RedisClient)

	// Controller
	movieController := controllers.NewMovieController(movieService)
	userController := controllers.NewUserController(userService)

	// Swagger route
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Register Routes
	routes.RegisterRoutes(e, movieController, userController)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	e.Logger.Fatal(e.Start(":" + port))
}
