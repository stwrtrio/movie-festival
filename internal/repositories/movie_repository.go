package repositories

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/stwrtrio/movie-festival/internal/models"
)

type MovieRepository interface {
	Create(ctx context.Context, movie *models.Movie) error
	Update(ctx context.Context, movie *models.Movie) error
	GetMostViewedMovie(ctx context.Context) (*models.Movie, error)
	GetMostViewedGenre(ctx context.Context) (string, int, error)
	GetAllMovies(ctx context.Context, limit, offset int) ([]models.Movie, error)
}

type movieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) MovieRepository {
	return &movieRepository{db}
}

func (r *movieRepository) Create(ctx context.Context, movie *models.Movie) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// create id
	movie.ID = uuid.NewString()

	// Insert movie
	query := `
        INSERT INTO movies (id, title, description, duration, watch_url, views) 
        VALUES (?, ?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, query, movie.ID, movie.Title, movie.Description, movie.Duration, movie.WatchURL, movie.Views)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert genres
	for _, genre := range movie.Genres {
		genreQuery := `
			INSERT INTO genres (name) VALUES (?)
			ON DUPLICATE KEY UPDATE name=name
		`
		res, err := tx.ExecContext(ctx, genreQuery, genre.Name)
		if err != nil {
			tx.Rollback()
			return err
		}
		genreID, _ := res.LastInsertId()

		// Link genre to movie
		linkQuery := `INSERT INTO movie_genres (movie_id, genre_id)
						VALUES (?, ?)
						ON DUPLICATE KEY UPDATE movie_id=movie_id, genre_id=genre_id`
		_, err = tx.ExecContext(ctx, linkQuery, movie.ID, genreID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Insert artists and associate with the movie
	for _, artist := range movie.Artists {
		var artistID string
		err := tx.QueryRowContext(ctx, "SELECT id FROM artists WHERE name = ?", artist.Name).Scan(&artistID)
		if err == sql.ErrNoRows {
			// Artist doesn't exist
			artistID = uuid.NewString()
			_, err = tx.ExecContext(ctx, "INSERT INTO artists (id, name) VALUES (?, ?)", artistID, artist.Name)
			if err != nil {
				tx.Rollback()
				return err
			}

		} else if err != nil {
			tx.Rollback()
			return err
		}

		// Associate artist with the movie
		_, err = tx.ExecContext(ctx,
			"INSERT INTO movie_artists (movie_id, artist_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE movie_id=movie_id, artist_id=artist_id",
			movie.ID, artistID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (r *movieRepository) Update(ctx context.Context, movie *models.Movie) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Update movie
	query := `
        UPDATE movies SET title=?, description=?, duration=?, watch_url=?, views=? 
        WHERE id=?`
	_, err = tx.ExecContext(ctx, query,
		movie.Title,
		movie.Description,
		movie.Duration,
		movie.WatchURL,
		movie.Views,
		movie.ID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert genres
	for _, genre := range movie.Genres {
		genreQuery := `
			INSERT INTO genres (name) VALUES (?)
			ON DUPLICATE KEY UPDATE name=name
		`
		res, err := tx.ExecContext(ctx, genreQuery, genre.Name)
		if err != nil {
			tx.Rollback()
			return err
		}
		genreID, _ := res.LastInsertId()

		// Link genre to movie
		linkQuery := `INSERT INTO movie_genres (movie_id, genre_id)
						VALUES (?, ?)
						ON DUPLICATE KEY UPDATE movie_id=movie_id, genre_id=genre_id`
		_, err = tx.ExecContext(ctx, linkQuery, movie.ID, genreID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Insert artists and associate with the movie
	for _, artist := range movie.Artists {
		var artistID string
		err := tx.QueryRowContext(ctx, "SELECT id FROM artists WHERE name = ?", artist.Name).Scan(&artistID)
		if err == sql.ErrNoRows {
			// Artist doesn't exist
			artistID = uuid.NewString()
			_, err = tx.ExecContext(ctx, "INSERT INTO artists (id, name) VALUES (?, ?)", artistID, artist.Name)
			if err != nil {
				tx.Rollback()
				return err
			}

		} else if err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.ExecContext(ctx,
			"INSERT INTO movie_artists (movie_id, artist_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE movie_id=movie_id, artist_id=artist_id",
			movie.ID, artistID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (r *movieRepository) GetMostViewedMovie(ctx context.Context) (*models.Movie, error) {
	query := `
		SELECT m.id, m.title, m.description, m.duration, m.watch_url, mv.view_count
		FROM movies m
		JOIN movie_views mv ON m.id = mv.movie_id
		ORDER BY mv.view_count DESC
		LIMIT 1
	`
	var movie models.Movie
	err := r.db.QueryRowContext(ctx, query).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Duration,
		&movie.WatchURL,
		&movie.Views)
	if err != nil {
		return nil, err
	}

	// Query to fetch genres for the movie
	genreQuery := `
		SELECT g.id, g.name
		FROM genres g
		JOIN movie_genres mg ON g.id = mg.genre_id
		WHERE mg.movie_id = ?
	`
	rows, err := r.db.QueryContext(ctx, genreQuery, movie.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect genres
	genres := []models.Genre{}
	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(&genre.ID, &genre.Name); err != nil {
			return nil, err
		}

		genres = append(genres, genre)
	}

	// Assign genres to the movie model
	movie.Genres = genres

	return &movie, nil
}

func (r *movieRepository) GetMostViewedGenre(ctx context.Context) (string, int, error) {
	query := `
		SELECT g.name, SUM(mv.view_count) AS total_views
		FROM movie_genres mg
		JOIN genres g ON mg.genre_id = g.id
		JOIN movie_views mv ON mg.movie_id = mv.movie_id
		GROUP BY g.id
		ORDER BY total_views DESC
		LIMIT 1
	`

	var genreName string
	var totalViews int
	err := r.db.QueryRowContext(ctx, query).Scan(&genreName, &totalViews)
	if err != nil {
		return "", 0, err
	}
	return genreName, totalViews, nil
}

func (r *movieRepository) GetAllMovies(ctx context.Context, limit, offset int) ([]models.Movie, error) {
	// we will set default limit if user not input the limit
	if limit < 1 {
		limit = 10
	}

	query := `
		SELECT m.id, m.title, m.description, m.duration, m.watch_url, m.views
		FROM movies m
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.Duration,
			&movie.WatchURL,
			&movie.Views,
		)
		if err != nil {
			return nil, err
		}

		// Fetch genres for the movie
		genreQuery := `
			SELECT g.id, g.name
			FROM genres g
			JOIN movie_genres mg ON g.id = mg.genre_id
			WHERE mg.movie_id = ?
		`
		genreRows, err := r.db.QueryContext(ctx, genreQuery, movie.ID)
		if err != nil {
			return nil, err
		}
		defer genreRows.Close()

		for genreRows.Next() {
			var genre models.Genre
			err := genreRows.Scan(&genre.ID, &genre.Name)
			if err != nil {
				return nil, err
			}
			movie.Genres = append(movie.Genres, genre)
		}

		movies = append(movies, movie)
	}

	return movies, nil
}
