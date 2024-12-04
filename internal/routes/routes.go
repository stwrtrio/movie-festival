package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/stwrtrio/movie-festival/internal/controllers"
	"github.com/stwrtrio/movie-festival/internal/middlewares"
)

func RegisterRoutes(e *echo.Echo, movieController *controllers.MovieController, userController *controllers.UserController) {

	// Public routes (no authentication required)
	e.POST("/api/user/register", userController.Register)
	e.POST("/api/user/login", userController.Login)

	e.POST("/api/movies/:id/view", movieController.TrackMovieView)
	e.GET("/api/movies", movieController.GetAllMovies)
	e.GET("/api/movies/search", movieController.SearchMovies)

	// Authenticated user routes
	userGroup := e.Group("/api/user")
	userGroup.Use(middlewares.AuthMiddleware)
	userGroup.POST("/logout", userController.Logout)
	userGroup.POST("/movies/:id/vote", movieController.VoteMovie)
	userGroup.POST("/movies/:id/unvote", movieController.UnvoteMovie)
	userGroup.GET("/votes", movieController.GetUserVotesController)

	// Admin routes
	adminGroup := e.Group("/api/admin")
	adminGroup.Use(middlewares.AdminAuthMiddleware)
	adminGroup.POST("/movie", movieController.CreateMovie)
	adminGroup.POST("/movie/:id", movieController.UpdateMovie)
	adminGroup.GET("/movies/most-viewed", movieController.GetMostViewedMovie)
	adminGroup.GET("/movies/most-viewed-genres", movieController.GetMostViewedGenre)
}
