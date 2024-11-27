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

	// Initialize Echo
	e := echo.New()
	e.Validator = &middlewares.CustomValidator{Validator: validator.New()}

	// Dependency Injection
	userRepo := repositories.NewMovieRepository(config.DB)
	userService := services.NewMovieService(userRepo)
	userController := controllers.NewMovieController(userService)

	// Register Routes
	routes.RegisterMovieRoutes(e, userController)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	e.Logger.Fatal(e.Start(":" + port))
}
