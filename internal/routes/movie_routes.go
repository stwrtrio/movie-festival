package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/stwrtrio/movie-festival/internal/controllers"
)

func RegisterMovieRoutes(e *echo.Echo, movieController *controllers.MovieController) {
	// Admin api
	e.POST("/api/admin/movies", movieController.CreateMovie)
	e.POST("/api/admin/movies/:id", movieController.UpdateMovie)
	e.GET("/api/admin/movies/most-viewed", movieController.GetMostViewedMovie)
	e.GET("/api/admin/movies/most-viewed-genre", movieController.GetMostViewedGenre)

	// Public api
	e.GET("/api/movies", movieController.GetAllMovies)
	e.GET("/api/movies/search", movieController.SearchMovies)
}
