package services

import (
	"context"

	"github.com/stwrtrio/movie-festival/internal/models"
	"github.com/stwrtrio/movie-festival/internal/repositories"
)

type MovieService interface {
	CreateMovie(ctx context.Context, movie *models.Movie) error
	UpdateMovie(ctx context.Context, movie *models.Movie) error
	GetMostViewedMovie(ctx context.Context) (*models.Movie, error)
	GetMostViewedGenre(ctx context.Context) (string, int, error)
}

type movieService struct {
	repo repositories.MovieRepository
}

func NewMovieService(repo repositories.MovieRepository) MovieService {
	return &movieService{repo}
}

func (s *movieService) CreateMovie(ctx context.Context, movie *models.Movie) error {
	return s.repo.Create(ctx, movie)
}

func (s *movieService) UpdateMovie(ctx context.Context, movie *models.Movie) error {
	return s.repo.Update(ctx, movie)
}

func (s *movieService) GetMostViewedMovie(ctx context.Context) (*models.Movie, error) {
	return s.repo.GetMostViewedMovie(ctx)
}

func (s *movieService) GetMostViewedGenre(ctx context.Context) (string, int, error) {
	return s.repo.GetMostViewedGenre(ctx)
}
