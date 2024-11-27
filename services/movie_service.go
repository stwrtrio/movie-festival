package services

import (
	"github.com/stwrtrio/movie-festival/models"
	"github.com/stwrtrio/movie-festival/repositories"
)

type MovieService interface {
	CreateMovie(movie *models.Movie) error
	UpdateMovie(movie *models.Movie) error
}

type movieService struct {
	repo repositories.MovieRepository
}

func NewMovieService(repo repositories.MovieRepository) MovieService {
	return &movieService{repo}
}

func (s *movieService) CreateMovie(movie *models.Movie) error {
	return s.repo.Create(movie)
}

func (s *movieService) UpdateMovie(movie *models.Movie) error {
	return s.repo.Update(movie)
}
