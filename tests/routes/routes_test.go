package routes_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stwrtrio/movie-festival/internal/controllers"
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

func TestMovieRoutes(t *testing.T) {
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
			"genres":      "Action,Thriller",
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
