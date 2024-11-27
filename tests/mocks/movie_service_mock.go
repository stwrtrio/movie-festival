package mocks

import (
	"github.com/stwrtrio/movie-festival/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockMovieService struct {
	mock.Mock
}

func (m *MockMovieService) CreateMovie(movie *models.Movie) error {
	args := m.Called(movie)
	return args.Error(0)
}

func (m *MockMovieService) UpdateMovie(movie *models.Movie) error {
	args := m.Called(movie)
	return args.Error(0)
}
