package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/stwrtrio/movie-festival/internal/models"
	"github.com/stwrtrio/movie-festival/internal/repositories"

	"github.com/go-redis/redis/v8"
)

type MovieService interface {
	CreateMovie(ctx context.Context, movie *models.Movie) error
	UpdateMovie(ctx context.Context, movie *models.Movie) error
	GetMostViewedMovie(ctx context.Context) (*models.Movie, error)
	GetMostViewedGenre(ctx context.Context) (string, int, error)
	GetAllMovies(ctx context.Context, limit, offset int) ([]models.Movie, error)
	GetAllMoviesFromCache(ctx context.Context, limit, offset int) ([]models.Movie, error)
	SearchMovies(ctx context.Context, query string, limit, offset int) ([]models.Movie, error)
	TrackMovieView(ctx context.Context, movieID string) error
}

type movieService struct {
	repo  repositories.MovieRepository
	redis redis.Cmdable
}

func NewMovieService(repo repositories.MovieRepository, redisClient redis.Cmdable) MovieService {
	return &movieService{repo: repo, redis: redisClient}
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

// GetAllMovies fetches movies from the database
func (s *movieService) GetAllMovies(ctx context.Context, limit, offset int) ([]models.Movie, error) {
	return s.repo.GetAllMovies(ctx, limit, offset)
}

// GetAllMoviesFromCache tries to fetch movies from Redis, and falls back to database if not found
func (s *movieService) GetAllMoviesFromCache(ctx context.Context, limit, offset int) ([]models.Movie, error) {
	cacheKey := fmt.Sprintf("movies:limit=%d:offset=%d", limit, offset)

	// Try to get movies from cache
	cacheData, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		// If found in cache, unmarshal and return
		var cachedMovies []models.Movie
		err := json.Unmarshal([]byte(cacheData), &cachedMovies)
		if err == nil {
			return cachedMovies, nil
		}
	}

	// If not found in cache, fetch from database
	movies, err := s.repo.GetAllMovies(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Cache the movies for future requests
	cacheByte, err := json.Marshal(movies)
	if err == nil {
		expiredAt, err := time.ParseDuration(os.Getenv("CACHE_DEFAULT_EXPIRATION"))
		if err != nil {
			return nil, err
		}

		s.redis.Set(ctx, cacheKey, string(cacheByte), expiredAt) // Store in Redis with expiration
	}

	return movies, nil
}

func (s *movieService) SearchMovies(ctx context.Context, query string, limit, offset int) ([]models.Movie, error) {
	return s.repo.SearchMovies(ctx, query, limit, offset)
}

func (s *movieService) TrackMovieView(ctx context.Context, movieID string) error {
	return s.repo.TrackMovieView(ctx, movieID)
}
