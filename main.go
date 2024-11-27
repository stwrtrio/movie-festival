package main

import (
	"log"

	"github.com/stwrtrio/movie-festival/config"
	"github.com/stwrtrio/movie-festival/controllers"
	"github.com/stwrtrio/movie-festival/middlewares"
	"github.com/stwrtrio/movie-festival/repositories"
	"github.com/stwrtrio/movie-festival/routes"
	"github.com/stwrtrio/movie-festival/services"

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

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
