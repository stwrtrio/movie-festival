package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stwrtrio/movie-festival/internal/controllers"
	"github.com/stwrtrio/movie-festival/internal/models"
	"github.com/stwrtrio/movie-festival/tests/mocks"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMovieService is a mock implementation of the MovieService interface.
type MockMovieService struct {
	mock.Mock
}

func NewValidator() *CustomValidator {
	v := validator.New()

	return &CustomValidator{Validator: v}
}

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

// TestCreateMovie tests the CreateMovie method of the MovieController.
func TestCreateMovie(t *testing.T) {
	// Create a mock service
	mockService := new(mocks.MockMovieService)
	controller := controllers.NewMovieController(mockService)

	// Prepare test input and expected output
	movieRequest := &models.CreateMovieRequest{
		Title:       "Test Movie",
		Description: "A great movie",
		Duration:    120,
		Genres:      "Action,Thriller",
		WatchURL:    "http://example.com/movie.mp4",
		Artists:     []string{"Actor A", "Actor B"},
	}

	movieModel := &models.Movie{
		Title:       movieRequest.Title,
		Description: movieRequest.Description,
		Duration:    movieRequest.Duration,
		Genres:      movieRequest.Genres,
		WatchURL:    movieRequest.WatchURL,
		Artists: []models.Artist{
			{Name: "Actor A"},
			{Name: "Actor B"},
		},
	}

	// Mock the service call
	mockService.On("CreateMovie", movieModel).Return(nil)

	e := echo.New()
	e.Validator = NewValidator()
	reqBody, _ := json.Marshal(movieRequest)
	req := httptest.NewRequest(http.MethodPost, "/api/admin/movies", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	res := rec.Result()
	defer res.Body.Close()

	// Call the CreateMovie method
	err := controller.CreateMovie(ctx)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	mockService.AssertExpectations(t)
}

// TestCreateMovie_InvalidRequest tests the CreateMovie method with invalid input.
func TestCreateMovie_InvalidRequest(t *testing.T) {
	// Create a mock service
	mockService := new(mocks.MockMovieService)

	// Create the controller with the mock service
	controller := controllers.NewMovieController(mockService)

	// Prepare invalid input (missing required fields)
	invalidRequest := map[string]interface{}{
		"Description": "A movie without a title",
	}
	reqBody, _ := json.Marshal(invalidRequest)

	// Create a new Echo context with a request and recorder
	e := echo.New()
	e.Validator = NewValidator() // Add validator
	req := httptest.NewRequest(http.MethodPost, "/api/admin/movies", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// Call the CreateMovie method
	err := controller.CreateMovie(ctx)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Key: 'CreateMovieRequest.Title' Error:Field validation for 'Title'")
	mockService.AssertExpectations(t)
}

// TestCreateMovie_ServiceError tests the CreateMovie method when the service returns an error.
func TestCreateMovie_ServiceError(t *testing.T) {
	// Create a mock service
	mockService := new(mocks.MockMovieService)

	// Create the controller with the mock service
	controller := controllers.NewMovieController(mockService)

	// Prepare test input and service error
	movieRequest := &models.CreateMovieRequest{
		Title:       "Test Movie",
		Description: "A great movie",
		Duration:    120,
		Genres:      "Action,Thriller",
		WatchURL:    "http://example.com/movie.mp4",
		Artists:     []string{"Actor A", "Actor B"},
	}
	movieModel := &models.Movie{
		Title:       movieRequest.Title,
		Description: movieRequest.Description,
		Duration:    movieRequest.Duration,
		Genres:      movieRequest.Genres,
		WatchURL:    movieRequest.WatchURL,
		Artists: []models.Artist{
			{Name: "Actor A"},
			{Name: "Actor B"},
		},
	}

	mockService.On("CreateMovie", movieModel).Return(errors.New("service error"))

	// Create a new Echo context with a request and recorder
	e := echo.New()
	e.Validator = NewValidator() // Add validator
	reqBody, _ := json.Marshal(movieRequest)
	req := httptest.NewRequest(http.MethodPost, "/api/admin/movies", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// Call the CreateMovie method
	err := controller.CreateMovie(ctx)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Failed to create movie")
	mockService.AssertExpectations(t)
}
