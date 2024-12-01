package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/stwrtrio/movie-festival/internal/controllers"
)

func RegisterRoutes(e *echo.Echo, movieController *controllers.MovieController, userController *controllers.UserController) {
	// Admin api
	e.POST("/api/admin/movies", movieController.CreateMovie)
	e.POST("/api/admin/movies/:id", movieController.UpdateMovie)
	e.GET("/api/admin/movies/most-viewed", movieController.GetMostViewedMovie)
	e.GET("/api/admin/movies/most-viewed-genres", movieController.GetMostViewedGenre)

	// Public api
	e.GET("/api/movies", movieController.GetAllMovies)
	e.GET("/api/movies/search", movieController.SearchMovies)
	e.POST("/api/movies/:id/view", movieController.TrackMovieView)

	// User api
	e.POST("/api/user/login", userController.Login)
}
