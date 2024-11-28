package mocks

import (
	"context"

	"github.com/stwrtrio/movie-festival/internal/models"

	"github.com/stretchr/testify/mock"
)

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

func (m *MockMovieRepository) GetAllMovies(ctx context.Context, limit, offset int) ([]models.Movie, error) {
	args := m.Called(ctx, limit, offset)
	if movie, ok := args.Get(0).([]models.Movie); ok {
		return movie, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMovieRepository) GetAllMoviesFromCache(ctx context.Context, limit, offset int) ([]models.Movie, error) {
	args := m.Called(ctx, limit, offset)
	if movie, ok := args.Get(0).([]models.Movie); ok {
		return movie, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMovieRepository) SearchMovies(ctx context.Context, query string, limit, offset int) ([]models.Movie, error) {
	args := m.Called(ctx, limit, offset)
	if movie, ok := args.Get(0).([]models.Movie); ok {
		return movie, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMovieRepository) TrackMovieView(ctx context.Context, movieID string) error {
	args := m.Called(ctx, movieID)
	return args.Error(0)
}
