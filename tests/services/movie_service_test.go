package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stwrtrio/movie-festival/internal/models"
	"github.com/stwrtrio/movie-festival/internal/services"
)

// MockMovieRepository is a mock implementation of the MovieRepository interface.
type MockMovieRepository struct {
	mock.Mock
}

func (m *MockMovieRepository) Create(ctx context.Context, movie *models.Movie) error {
	args := m.Called(movie)
	return args.Error(0)
}

func (m *MockMovieRepository) Update(ctx context.Context, movie *models.Movie) error {
	args := m.Called(movie)
	return args.Error(0)
}

func (m *MockMovieRepository) GetMostViewedMovie(ctx context.Context) (*models.Movie, error) {
	args := m.Called(ctx)
	if movie, ok := args.Get(0).(*models.Movie); ok {
		return movie, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMovieRepository) GetMostViewedGenre(ctx context.Context) (string, int, error) {
	args := m.Called(ctx)
	return args.String(0), args.Int(0), args.Error(1)
}

// TestCreateMovie tests the CreateMovie method of the MovieService.
func TestCreateMovie(t *testing.T) {
	mockRepo := new(MockMovieRepository)
	service := services.NewMovieService(mockRepo)

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
	mockRepo := new(MockMovieRepository)
	service := services.NewMovieService(mockRepo)

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
	mockRepo := new(MockMovieRepository)
	service := services.NewMovieService(mockRepo)

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
	mockRepo := new(MockMovieRepository)
	service := services.NewMovieService(mockRepo)

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
