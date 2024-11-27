package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/stwrtrio/movie-festival/controllers"
)

func RegisterMovieRoutes(e *echo.Echo, movieController *controllers.MovieController) {
	e.POST("/api/admin/movies", movieController.CreateMovie)
	e.POST("/api/admin/movies/:id", movieController.UpdateMovie)
}
