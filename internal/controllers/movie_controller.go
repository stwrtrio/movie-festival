package controllers

import (
	"net/http"

	"github.com/stwrtrio/movie-festival/internal/models"
	"github.com/stwrtrio/movie-festival/internal/services"
	"github.com/stwrtrio/movie-festival/internal/utils"

	"github.com/labstack/echo/v4"
)

type MovieController struct {
	service services.MovieService
}

func NewMovieController(service services.MovieService) *MovieController {
	return &MovieController{service}
}

func (c *MovieController) CreateMovie(ctx echo.Context) error {
	req := new(models.CreateMovieRequest)
	if err := ctx.Bind(req); err != nil {
		return utils.FailResponse(ctx, "Invalid request body", http.StatusBadRequest)
	}
	if err := ctx.Validate(req); err != nil {
		return utils.FailResponse(ctx, err.Error(), http.StatusBadRequest)
	}

	if len(req.Artists) == 0 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "At least one artist is required"})
	}

	// Convert request data to a Movie model
	artists := make([]models.Artist, 0)
	for _, artistName := range req.Artists {
		artists = append(artists, models.Artist{Name: string(artistName)})
	}

	movie := &models.Movie{
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		Genres:      req.Genres,
		WatchURL:    req.WatchURL,
		Artists:     artists,
	}

	if err := c.service.CreateMovie(movie); err != nil {
		return utils.FailResponse(ctx, "Failed to create movie", http.StatusInternalServerError)
	}

	return utils.SuccessResponse(ctx, http.StatusCreated, "Movie created successfully", nil)
}

func (c *MovieController) UpdateMovie(ctx echo.Context) error {
	movieID := ctx.Param("id")
	if movieID == "" {
		return utils.FailResponse(ctx, "Invalid request body", http.StatusBadRequest)
	}

	req := new(models.CreateMovieRequest)
	if err := ctx.Bind(req); err != nil {
		return utils.FailResponse(ctx, "Invalid request body", http.StatusBadRequest)
	}

	// Convert request data to a Movie model
	artists := make([]models.Artist, 0)
	for _, artistName := range req.Artists {
		artists = append(artists, models.Artist{Name: string(artistName)})
	}

	movie := &models.Movie{
		ID:          movieID,
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		Genres:      req.Genres,
		WatchURL:    req.WatchURL,
		Artists:     artists,
	}

	if err := c.service.UpdateMovie(movie); err != nil {
		return utils.FailResponse(ctx, err.Error(), http.StatusInternalServerError)
	}

	return utils.SuccessResponse(ctx, http.StatusCreated, "Movie updated successfully", nil)
}
