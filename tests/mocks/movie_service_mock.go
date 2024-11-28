package mocks

import (
	"context"

	"github.com/stwrtrio/movie-festival/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockMovieService is a mock implementation of the MovieService interface.
type MockMovieService struct {
	mock.Mock
}

func (m *MockMovieService) CreateMovie(ctx context.Context, movie *models.Movie) error {
	args := m.Called(movie)
	return args.Error(0)
}

func (m *MockMovieService) UpdateMovie(ctx context.Context, movie *models.Movie) error {
	args := m.Called(movie)
	return args.Error(0)
}

func (m *MockMovieService) GetMostViewedMovie(ctx context.Context) (*models.Movie, error) {
	args := m.Called(ctx)
	if movie, ok := args.Get(0).(*models.Movie); ok {
		return movie, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMovieService) GetMostViewedGenre(ctx context.Context) (string, int, error) {
	args := m.Called(ctx)
	return args.String(0), args.Int(0), args.Error(1)
}

func (m *MockMovieService) GetAllMovies(ctx context.Context, limit, offset int) ([]models.Movie, error) {
	args := m.Called(ctx)
	if movie, ok := args.Get(0).([]models.Movie); ok {
		return movie, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMovieService) GetAllMoviesFromCache(ctx context.Context, limit, offset int) ([]models.Movie, error) {
	args := m.Called(ctx)
	if movie, ok := args.Get(0).([]models.Movie); ok {
		return movie, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMovieService) SearchMovies(ctx context.Context, query string, limit, offset int) ([]models.Movie, error) {
	args := m.Called(ctx)
	if movie, ok := args.Get(0).([]models.Movie); ok {
		return movie, args.Error(1)
	}
	return nil, args.Error(1)
}
