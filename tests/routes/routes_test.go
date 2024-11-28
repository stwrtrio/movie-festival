package routes_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stwrtrio/movie-festival/internal/controllers"
	"github.com/stwrtrio/movie-festival/internal/models"
	"github.com/stwrtrio/movie-festival/internal/routes"
	"github.com/stwrtrio/movie-festival/tests/mocks"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// CustomValidator for Echo
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func TestCreateMovieRoutes(t *testing.T) {
	// Initialize Echo and mock service
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	mockService := new(mocks.MockMovieService)

	// Set up routes with mock controller
	movieController := controllers.NewMovieController(mockService)
	routes.RegisterMovieRoutes(e, movieController)

	t.Run("Test Create Movie Route", func(t *testing.T) {
		// Create a sample request payload
		requestPayload := map[string]interface{}{
			"title":       "Test Movie",
			"description": "A test movie",
			"duration":    120,
			"genres":      []string{"Action", "Thriller"},
			"artists":     []string{"Actor 1", "Actor 2"},
			"watch_url":   "http://example.com/test.mp4",
		}
		body, _ := json.Marshal(requestPayload)

		// Create a new HTTP POST request
		req := httptest.NewRequest(http.MethodPost, "/api/admin/movies", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		// Define expected behavior for the mock service
		mockService.On("CreateMovie", mock.Anything).Return(nil)

		// Simulate the request
		e.ServeHTTP(rec, req)

		// Assertions
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "Movie created successfully")

		// Ensure the mock was called as expected
		mockService.AssertExpectations(t)
	})

	t.Run("Test Invalid Payload for Create Movie", func(t *testing.T) {
		// Create an invalid payload
		invalidPayload := map[string]interface{}{
			"title": 123, // Invalid type for title
		}
		body, _ := json.Marshal(invalidPayload)

		// Create a new HTTP POST request
		req := httptest.NewRequest(http.MethodPost, "/api/admin/movies", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		// Simulate the request
		e.ServeHTTP(rec, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid request body")
	})
}

func TestGetMostViewedMovieRoute(t *testing.T) {
	// Create a mock service
	mockService := new(mocks.MockMovieService)
	controller := controllers.NewMovieController(mockService)

	// Prepare mock data
	mostViewedMovie := &models.Movie{
		ID:          uuid.NewString(),
		Title:       "Most Viewed Movie",
		Description: "This is the most viewed movie.",
		Duration:    120,
		WatchURL:    "http://example.com/movie.mp4",
		Genres:      []models.Genre{{Name: "Action"}, {Name: "Thriller"}},
		Views:       500,
	}
	mockService.On("GetMostViewedMovie", mock.Anything).Return(mostViewedMovie, nil)

	// Create an Echo instance and register the route
	e := echo.New()

	// Define the route in the test
	e.GET("/api/movies/most-viewed-movie", controller.GetMostViewedMovie)

	// Run subtest for successful response
	t.Run("Success", func(t *testing.T) {
		// Simulate a GET request to the route
		req := httptest.NewRequest(http.MethodGet, "/api/movies/most-viewed-movie", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		// Assertions
		assert.Equal(t, http.StatusOK, rec.Code)

		// Parse the response body
		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check if the response contains the correct movie data
		data, ok := response["data"].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected 'data' to be an object, but got %T", response["data"])
		}

		// Validate the movie data
		assert.Equal(t, mostViewedMovie.Title, data["title"])
		assert.Equal(t, mostViewedMovie.Description, data["description"])
		assert.Equal(t, mostViewedMovie.Duration, int(data["duration"].(float64)))
		assert.Equal(t, mostViewedMovie.WatchURL, data["watch_url"])
		assert.Equal(t, mostViewedMovie.Views, int(data["views"].(float64)))

		// Validate genres
		genres, ok := data["genres"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, genres, len(mostViewedMovie.Genres))

		// Validate the first genre
		genre := genres[0].(map[string]interface{})
		assert.Equal(t, mostViewedMovie.Genres[0].Name, genre["name"])
	})

	// Run subtest for error response
	t.Run("Error", func(t *testing.T) {
		// Mock the service to return an error
		mockService.On("GetMostViewedMovie", mock.Anything).Return(nil, errors.New("failed to get most viewed movie"))

		// Simulate a GET request to the route
		req := httptest.NewRequest(http.MethodGet, "/api/movies/most-viewed-movies", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		// Assertions
		assert.Equal(t, http.StatusNotFound, rec.Code)

		// Parse the response body
		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check if the error message is returned
		errorMessage, ok := response["message"].(string)
		if !ok {
			t.Fatalf("Expected 'message' to be a string, but got %T", response["message"])
		}

		// Validate the error response
		assert.Equal(t, "Not Found", errorMessage)
	})

	mockService.AssertExpectations(t)
}
