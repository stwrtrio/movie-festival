package services_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stwrtrio/movie-festival/internal/models"
	"github.com/stwrtrio/movie-festival/internal/services"
	"github.com/stwrtrio/movie-festival/tests/mocks"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRedisClient struct {
	mock.Mock
}

// TestCreateMovie tests the CreateMovie method of the MovieService.
func TestCreateMovie(t *testing.T) {
	// Create a mock repository
	mockRepo := new(mocks.MockMovieRepository)
	// Mock Redis client
	mockRedisClient, _ := redismock.NewClientMock()

	service := services.NewMovieService(mockRepo, mockRedisClient)

	movie := &models.Movie{
		Title:       "Test Movie",
		Description: "A test movie",
		Duration:    120,
		Genres: []models.Genre{
			{Name: "Action"},
		},
		WatchURL: "http://example.com/test.mp4",
	}

	// Mock repository behavior
	mockRepo.On("Create", movie).Return(nil)

	// Call the service method
	err := service.CreateMovie(context.Background(), movie)

	// Assertions
	assert.Nil(t, err)
	mockRepo.AssertExpectations(t)
}

// TestCreateMovie_Fail tests the CreateMovie method when the repository returns an error.
func TestCreateMovie_Fail(t *testing.T) {
	mockRepo := new(mocks.MockMovieRepository)
	// Mock Redis client
	mockRedisClient, _ := redismock.NewClientMock()
	service := services.NewMovieService(mockRepo, mockRedisClient)

	movie := &models.Movie{
		Title:       "Test Movie",
		Description: "A test movie",
		Duration:    120,
		Genres: []models.Genre{
			{Name: "Action"},
		},
		WatchURL: "http://example.com/test.mp4",
	}

	// Mock repository behavior to return an error
	mockRepo.On("Create", movie).Return(errors.New("repository error"))

	// Call the service method
	err := service.CreateMovie(context.Background(), movie)

	// Assertions
	assert.NotNil(t, err)
	assert.Equal(t, "repository error", err.Error())
	mockRepo.AssertExpectations(t)
}

// TestUpdateMovie tests the UpdateMovie method of the MovieService.
func TestUpdateMovie(t *testing.T) {
	mockRepo := new(mocks.MockMovieRepository)
	// Mock Redis client
	mockRedisClient, _ := redismock.NewClientMock()
	service := services.NewMovieService(mockRepo, mockRedisClient)

	movie := &models.Movie{
		ID:          "123",
		Title:       "Updated Movie",
		Description: "An updated movie",
		Duration:    140,
		Genres: []models.Genre{
			{Name: "Drama"},
		},
		WatchURL: "http://example.com/updated.mp4",
	}

	// Mock repository behavior
	mockRepo.On("Update", movie).Return(nil)

	// Call the service method
	err := service.UpdateMovie(context.Background(), movie)

	// Assertions
	assert.Nil(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUpdateMovie_Fail tests the UpdateMovie method when the repository returns an error.
func TestUpdateMovie_Fail(t *testing.T) {
	mockRepo := new(mocks.MockMovieRepository)
	mockRedisClient, _ := redismock.NewClientMock()
	service := services.NewMovieService(mockRepo, mockRedisClient)

	movie := &models.Movie{
		ID:          "123",
		Title:       "Updated Movie",
		Description: "An updated movie",
		Duration:    140,
		Genres: []models.Genre{
			{Name: "Drama"},
		},
		WatchURL: "http://example.com/updated.mp4",
	}

	// Mock repository behavior to return an error
	mockRepo.On("Update", movie).Return(errors.New("update failed"))

	// Call the service method
	err := service.UpdateMovie(context.Background(), movie)

	// Assertions
	assert.NotNil(t, err)
	assert.Equal(t, "update failed", err.Error())
	mockRepo.AssertExpectations(t)
}

// TestGetAllMovies tests the GetAllMovies method of the MovieService.
func TestGetAllMovies(t *testing.T) {
	// Create mock repository
	mockRepo := new(mocks.MockMovieRepository)

	// Prepare test data (movies)
	movies := []models.Movie{
		{ID: "1", Title: "Movie 1", Description: "Description 1"},
		{ID: "2", Title: "Movie 2", Description: "Description 2"},
	}

	// Mock repository call to return movie data
	mockRepo.On("GetAllMovies", mock.Anything, 10, 0).Return(movies, nil).Once()

	service := services.NewMovieService(mockRepo, nil)
	result, err := service.GetAllMovies(context.Background(), 10, 0)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "Movie 1", result[0].Title)
	assert.Equal(t, "Movie 2", result[1].Title)

	mockRepo.AssertExpectations(t)
}

func TestGetAllMoviesFromCache(t *testing.T) {
	// Create mock repository and Redis client
	mockRepo := new(mocks.MockMovieRepository)
	mockRedisClient, mockRedis := redismock.NewClientMock()

	// Prepare test data (movies)
	movies := []models.Movie{
		{ID: "1", Title: "Movie 1", Description: "Description 1"},
		{ID: "2", Title: "Movie 2", Description: "Description 2"},
	}

	// Cache key for GetAllMoviesFromCache
	cacheKey := "movies:limit=10:offset=0"
	mockRedis.ExpectGet(cacheKey).RedisNil()

	mockRepo.On("GetAllMovies", mock.Anything, 10, 0).Return(movies, nil).Once()

	// Set the environment variable for cache expiration
	os.Setenv("CACHE_DEFAULT_EXPIRATION", "1m")
	cacheData, err := json.Marshal(movies)
	if err != nil {
		t.Fatal("Failed to marshal movie data", err)
	}
	mockRedis.ExpectSet(cacheKey, string(cacheData), time.Minute).SetVal("OK") // Cache expiration 1 minute

	service := services.NewMovieService(mockRepo, mockRedisClient)
	result, err := service.GetAllMoviesFromCache(context.Background(), 10, 0)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "Movie 1", result[0].Title)
	assert.Equal(t, "Movie 2", result[1].Title)

	mockRepo.AssertExpectations(t)
}
