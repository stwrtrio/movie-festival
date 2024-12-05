package controllers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/stwrtrio/movie-festival/internal/middlewares"
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
	cx := ctx.Request().Context()
	req := new(models.CreateMovieRequest)
	if err := ctx.Bind(req); err != nil {
		return utils.FailResponse(ctx, http.StatusBadRequest, "Invalid request body")
	}
	if err := ctx.Validate(req); err != nil {
		log.Printf("Validation error: %v", err)
		return utils.FailResponse(ctx, http.StatusBadRequest, err.Error())
	}

	// Convert request data to a Genre model
	genres := make([]models.Genre, 0)
	for _, genreName := range req.Genres {
		genres = append(genres, models.Genre{Name: string(genreName)})
	}

	// Convert request data to a Artist model
	artists := make([]models.Artist, 0)
	for _, artistName := range req.Artists {
		artists = append(artists, models.Artist{Name: string(artistName)})
	}

	movie := &models.Movie{
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		Genres:      genres,
		WatchURL:    req.WatchURL,
		Artists:     artists,
	}

	if err := c.service.CreateMovie(cx, movie); err != nil {
		if err == errors.New("service CreateMovie err: movie doesn't have artist") {
			return utils.SuccessResponse(ctx, http.StatusCreated, "Failed to create movie: movie doesn't have artist", nil)
		}
		return utils.FailResponse(ctx, http.StatusInternalServerError, "Failed to create movie")
	}

	return utils.SuccessResponse(ctx, http.StatusCreated, "Movie created successfully", nil)
}

func (c *MovieController) UpdateMovie(ctx echo.Context) error {
	cx := ctx.Request().Context()
	movieID := ctx.Param("id")
	if movieID == "" {
		return utils.FailResponse(ctx, http.StatusBadRequest, "Invalid request body")
	}

	req := new(models.CreateMovieRequest)
	if err := ctx.Bind(req); err != nil {
		return utils.FailResponse(ctx, http.StatusBadRequest, "Invalid request body")
	}

	// Convert request data to a Movie model
	genres := make([]models.Genre, 0)
	for _, genreName := range req.Genres {
		genres = append(genres, models.Genre{Name: string(genreName)})
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
		Genres:      genres,
		WatchURL:    req.WatchURL,
		Artists:     artists,
	}

	err := c.service.UpdateMovie(cx, movie)
	if err != nil {
		if err == sql.ErrNoRows {
			return utils.FailResponse(ctx, http.StatusBadRequest, "movie is not exists")
		}

		return utils.FailResponse(ctx, http.StatusInternalServerError, "Unexpected error occurred. Please contact support")
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "Movie updated successfully", nil)
}

func (c *MovieController) GetMostViewedMovie(ctx echo.Context) error {
	cx := ctx.Request().Context()
	movie, err := c.service.GetMostViewedMovie(cx)
	if err != nil {
		return utils.FailResponse(ctx, http.StatusInternalServerError, "Unexpected error occurred. Please contact support")
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "", movie)
}

func (c *MovieController) GetMostViewedGenre(ctx echo.Context) error {
	// Retrieve pagination and sorting parameters from the query string
	pageStr := ctx.QueryParam("page")
	pageSizeStr := ctx.QueryParam("page_size")
	sortOrder := ctx.QueryParam("sort_order")

	// Default values if parameters are not provided
	if pageStr == "" {
		pageStr = "1" // default to page 1
	}
	if pageSizeStr == "" {
		pageSizeStr = "10" // default to 10 items per page
	}
	if sortOrder == "" {
		sortOrder = "DESC" // default to descending order
	}

	// Convert the parameters to integers
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		return utils.FailResponse(ctx, http.StatusBadRequest, "Invalid page number")
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		return utils.FailResponse(ctx, http.StatusBadRequest, "Invalid page size")
	}

	// Call the service to get the most viewed genres
	genreViews, err := c.service.GetMostViewedGenre(ctx.Request().Context(), page, pageSize, sortOrder)
	if err != nil {
		return utils.FailResponse(ctx, http.StatusInternalServerError, "Unexpected error occurred. Please contact support")
	}

	// Return the result
	return utils.SuccessResponse(ctx, http.StatusOK, "", genreViews)
}

// GetAllMovies handles GET requests to fetch all movies with pagination
func (c *MovieController) GetAllMovies(ctx echo.Context) error {
	cx := ctx.Request().Context()
	// Get pagination parameters
	limit := 10 // default limit
	offset := 0 // default offset

	if l := ctx.QueryParam("limit"); l != "" {
		// Convert the limit from string to int
		if parsedLimit, err := strconv.Atoi(l); err == nil {
			limit = parsedLimit
		}
	}

	if o := ctx.QueryParam("offset"); o != "" {
		// Convert the offset from string to int
		if parsedOffset, err := strconv.Atoi(o); err == nil {
			offset = parsedOffset
		}
	}

	// Check if the `use-cache` flag is set
	useCache := ctx.QueryParam("use-cache")

	var movies []models.Movie
	var err error

	// If `use-cache` is "true" or "1", try to fetch data from Redis
	if useCache == "true" || useCache == "1" {
		movies, err = c.service.GetAllMoviesFromCache(cx, limit, offset)
	} else {
		// Otherwise, fetch from the database
		movies, err = c.service.GetAllMovies(cx, limit, offset)
	}

	if err != nil {
		return utils.FailResponse(ctx, http.StatusInternalServerError, "Unexpected error occurred. Please contact support")
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "", movies)
}

func (c *MovieController) SearchMovies(ctx echo.Context) error {
	query := ctx.QueryParam("query")
	limit, err := strconv.Atoi(ctx.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 10 // default limit
	}
	offset, err := strconv.Atoi(ctx.QueryParam("offset"))
	if err != nil || offset < 0 {
		offset = 0 // default offset
	}

	movies, err := c.service.SearchMovies(ctx.Request().Context(), query, limit, offset)
	if err != nil {
		return utils.FailResponse(ctx, http.StatusInternalServerError, "Unexpected error occurred. Please contact support")
	}

	if len(movies) < 1 {
		movies = []models.Movie{}
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "", movies)
}

func (c *MovieController) TrackMovieView(ctx echo.Context) error {
	movieID := ctx.Param("id")
	err := c.service.TrackMovieView(ctx.Request().Context(), movieID)
	if err != nil {
		return utils.FailResponse(ctx, http.StatusInternalServerError, "Unexpected error occurred. Please contact support")
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "Viewership tracked successfully", nil)
}

func (c *MovieController) VoteMovie(ctx echo.Context) error {
	cx := ctx.Request().Context()
	// Get user claims from context
	claims, ok := middlewares.GetUserFromContext(ctx)
	if !ok {
		return utils.FailResponse(ctx, http.StatusUnauthorized, "User not authenticated")
	}
	userID := claims.UserID

	// Get movie ID from the URL
	movieID := ctx.Param("id")
	if movieID == "" {
		return utils.FailResponse(ctx, http.StatusBadRequest, "Movie ID is required")
	}

	// Call the service to vote for the movie
	err := c.service.VoteMovie(cx, userID, movieID)
	if err != nil {
		if err == errors.New("you have already voted for this movie") {
			return utils.FailResponse(ctx, http.StatusBadRequest, err.Error())
		}
		return utils.FailResponse(ctx, http.StatusInternalServerError, "Failed to vote for movie")
	}

	// Success response
	return utils.SuccessResponse(ctx, http.StatusOK, "Movie voted successfully", nil)
}

func (c *MovieController) UnvoteMovie(ctx echo.Context) error {
	cx := ctx.Request().Context()

	// Get user claims from context
	claims, ok := middlewares.GetUserFromContext(ctx)
	if !ok {
		return utils.FailResponse(ctx, http.StatusUnauthorized, "User not authenticated")
	}
	userID := claims.UserID

	// Extract the movie ID from the request parameters
	movieID := ctx.Param("id")
	if movieID == "" {
		return utils.FailResponse(ctx, http.StatusBadRequest, "Movie ID is required")
	}

	// Call the service to unvote the movie
	err := c.service.UnvoteMovie(cx, userID, movieID)
	if err != nil {
		if err == errors.New("you haven't voted for this movie yet") {
			return utils.FailResponse(ctx, http.StatusBadRequest, err.Error())
		}
		return utils.FailResponse(ctx, http.StatusInternalServerError, "Failed to unvote for movie")
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "Movie unvoted successfully", nil)
}

// GetUserVotesController handles fetching the list of voted movies for a user.
func (c *MovieController) GetUserVotesController(ctx echo.Context) error {
	cx := ctx.Request().Context()

	// Get user claims from context
	claims, ok := middlewares.GetUserFromContext(ctx)
	if !ok {
		return utils.FailResponse(ctx, http.StatusUnauthorized, "User not authenticated")
	}
	userID := claims.UserID

	// Call the service to get the voted movies
	votedMovies, err := c.service.GetUserVotedMovies(cx, userID)
	if err != nil {
		return utils.FailResponse(ctx, http.StatusInternalServerError, "Failed to fetch voted movies")
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "Voted movies retrieved successfully", votedMovies)
}

func (c *MovieController) GetMostVotedMovie(ctx echo.Context) error {
	cx := ctx.Request().Context()

	votedMovies, err := c.service.GetMostVotedMovie(cx)
	if err != nil {
		return utils.FailResponse(ctx, http.StatusInternalServerError, "Failed to fetch most voted movies")
	}

	return utils.SuccessResponse(ctx, http.StatusOK, "Voted movies retrieved successfully", votedMovies)
}
