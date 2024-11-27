package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stwrtrio/movie-festival/config"
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

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
