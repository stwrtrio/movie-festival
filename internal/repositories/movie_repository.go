package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/stwrtrio/movie-festival/internal/models"
)

type MovieRepository interface {
	Create(ctx context.Context, movie *models.Movie) error
	Update(ctx context.Context, movie *models.Movie) error
	GetMostViewedMovie(ctx context.Context) (*models.Movie, error)
	GetMostViewedGenre(ctx context.Context, page int, pageSize int, sortOrder string) ([]models.GenreView, error)
	GetAllMovies(ctx context.Context, limit, offset int) ([]models.Movie, error)
	SearchMovies(ctx context.Context, query string, limit, offset int) ([]models.Movie, error)
	TrackMovieView(ctx context.Context, movieID string) error
	FindMovieByID(ctx context.Context, movieID string) (models.Movie, error)
	FindGenreByMovieID(ctx context.Context, movieID string) (models.Genre, error)
	FindArtistByMovieID(ctx context.Context, movieID string) (models.Artist, error)
	GetVoteByUserAndMovie(ctx context.Context, userID, movieID string) (*models.Vote, error)
	CreateVote(ctx context.Context, userID, movieID string) error
}

type movieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) MovieRepository {
	return &movieRepository{db}
}

func (r *movieRepository) FindMovieByID(ctx context.Context, movieID string) (models.Movie, error) {
	var movie models.Movie
	query := `SELECT id, title, description, duration, watch_url FROM movies WHERE id = ?`
	err := r.db.QueryRowContext(ctx, query, movieID).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Duration,
		&movie.WatchURL)
	if err != nil {
		return movie, err
	}

	return movie, nil
}

func (r *movieRepository) FindGenreByMovieID(ctx context.Context, movieID string) (models.Genre, error) {
	var genre models.Genre
	query := `SELECT g.id, g.name
				FROM genres g
				JOIN movie_genres mg ON g.id = mg.genre_id
				WHERE mg.movie_id = ?`

	err := r.db.QueryRowContext(ctx, query).Scan(
		movieID,
		&genre.ID,
		&genre.Name)
	if err != nil {
		return genre, err
	}

	return genre, nil
}

func (r *movieRepository) FindArtistByMovieID(ctx context.Context, movieID string) (models.Artist, error) {
	var artist models.Artist
	query := `SELECT a.id, a.name
				FROM artists a
				JOIN movie_artists ma ON a.id = ma.artist_id
				WHERE ma.movie_id = ?`

	err := r.db.QueryRowContext(ctx, query).Scan(
		movieID,
		&artist.ID,
		&artist.Name)
	if err != nil {
		return artist, err
	}

	return artist, nil
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

	// Insert movie
	query := `
        INSERT INTO movies (id, title, description, duration, watch_url, views) 
        VALUES (?, ?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, query, movie.ID, movie.Title, movie.Description, movie.Duration, movie.WatchURL, movie.Views)
	if err != nil {
		tx.Rollback()
		log.Printf("Error insert movie: %v", err)
		return err
	}

	// Insert genres
	for _, genre := range movie.Genres {
		// create genre
		genreID, tx, err := r.createMovieGenre(tx, genre)
		if err != nil {
			tx.Rollback()
			log.Printf("Error creating movie genre: %v", err)
			return err
		}

		// Link genre to movie
		linkQuery := `INSERT INTO movie_genres (movie_id, genre_id)
						VALUES (?, ?)
						ON DUPLICATE KEY UPDATE movie_id=movie_id, genre_id=genre_id`
		_, err = tx.ExecContext(ctx, linkQuery, movie.ID, genreID)
		if err != nil {
			tx.Rollback()
			log.Printf("Error link movie genres: %v", err)
			return err
		}
	}

	// Insert artists and associate with the movie
	for _, artist := range movie.Artists {
		// create artist
		artistID, tx, err := r.createArtist(tx, artist)
		if err != nil {
			tx.Rollback()
			log.Printf("Error creating artist: %v", err)
			return err
		}

		_, err = tx.ExecContext(ctx,
			"INSERT INTO movie_artists (movie_id, artist_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE movie_id=movie_id, artist_id=artist_id",
			movie.ID, artistID)
		if err != nil {
			tx.Rollback()
			log.Printf("Error link movie artists: %v", err)
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

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Update movie details (e.g., title)
	_, err = tx.ExecContext(ctx,
		"UPDATE movies SET title = ?, description = ?, duration = ?, watch_url = ? WHERE id = ?",
		movie.Title, movie.Description, movie.Duration, movie.WatchURL, movie.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 2. Clear existing genres and update
	_, err = tx.ExecContext(ctx, "DELETE FROM movie_genres WHERE movie_id = ?", movie.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, genre := range movie.Genres {
		// create genre
		genreID, tx, err := r.createMovieGenre(tx, genre)
		if err != nil {
			tx.Rollback()
			return err
		}

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

	// 3. Update artists
	_, err = tx.ExecContext(ctx, "DELETE FROM movie_artists WHERE movie_id = ?", movie.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, artist := range movie.Artists {
		// create artist
		artistID, tx, err := r.createArtist(tx, artist)
		if err != nil {
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
		return err
	}

	return nil
}

func (r *movieRepository) createMovieGenre(tx *sql.Tx, genre models.Genre) (int64, *sql.Tx, error) {
	var genreID int64
	err := tx.QueryRow(`SELECT id FROM genres WHERE name = ?`, genre.Name).Scan(&genreID)

	if err == sql.ErrNoRows {
		// insert
		query := `INSERT INTO genres (name) VALUES (?)`
		res, err := tx.Exec(query, genre.Name)
		if err != nil {
			tx.Rollback()
			return genreID, tx, err
		}

		// get the generated id
		genreID, err = res.LastInsertId()

		if err != nil {
			tx.Rollback()
			return genreID, tx, err
		}
	} else if err != nil {
		tx.Rollback()
		return genreID, tx, err
	}

	return genreID, tx, nil
}

func (r *movieRepository) createArtist(tx *sql.Tx, artist models.Artist) (string, *sql.Tx, error) {
	var artistID string
	err := tx.QueryRow("SELECT id FROM artists WHERE name = ?", artist.Name).Scan(&artistID)
	if err == sql.ErrNoRows {
		// Artist doesn't exist
		_, err = tx.Exec("INSERT INTO artists (id, name) VALUES (?, ?)", artist.ID, artist.Name)
		if err != nil {
			tx.Rollback()
			return artistID, tx, err
		}

		artistID = artist.ID

	} else if err != nil {
		tx.Rollback()
		return artistID, tx, err
	}

	return artistID, tx, nil
}

func (r *movieRepository) GetMostViewedMovie(ctx context.Context) (*models.Movie, error) {
	query := `
		SELECT m.id, m.title, m.description, m.duration, m.watch_url, mv.view_count, m.created_at, m.updated_at
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
		&movie.Views,
		&movie.CreatedAt,
		&movie.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Fetch genres for the movie
	genreRows, err := r.getGenresByMovieID(ctx, movie.ID)
	if err != nil {
		log.Fatalf("Error fetching genres: %v", err)
	}
	movie.Genres = genreRows

	return &movie, nil
}

func (r *movieRepository) GetMostViewedGenre(ctx context.Context, page int, pageSize int, sortOrder string) ([]models.GenreView, error) {
	// Calculate the offset for pagination
	offset := (page - 1) * pageSize

	// Updated query with pagination and sorting by total_views
	query := fmt.Sprintf(`
		SELECT g.name, SUM(mv.view_count) AS total_views
		FROM movie_genres mg
		JOIN genres g ON mg.genre_id = g.id
		JOIN movie_views mv ON mg.movie_id = mv.movie_id
		GROUP BY g.id
		ORDER BY total_views %s, g.name
		LIMIT ? OFFSET ?
	`, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.GenreView
	for rows.Next() {
		var genreView models.GenreView
		err := rows.Scan(&genreView.Name, &genreView.ViewCount)
		if err != nil {
			return nil, err
		}
		result = append(result, genreView)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *movieRepository) GetAllMovies(ctx context.Context, limit, offset int) ([]models.Movie, error) {
	// we will set default limit if user not input the limit
	if limit < 1 {
		limit = 10
	}

	query := `
		SELECT m.id, m.title, m.description, m.duration, m.watch_url, m.created_at, m.updated_at
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
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Fetch genres for the movie
		genreRows, err := r.getGenresByMovieID(ctx, movie.ID)
		if err != nil {
			log.Fatalf("Error fetching genres: %v", err)
		}
		movie.Genres = genreRows
		movies = append(movies, movie)
	}

	return movies, nil
}

func (r *movieRepository) SearchMovies(ctx context.Context, query string, limit, offset int) ([]models.Movie, error) {
	query = "%" + query + "%"
	queryString := `
		SELECT DISTINCT m.id, m.title, m.description, m.duration, m.watch_url, m.created_at, m.updated_at 
		FROM movies m
		LEFT JOIN movie_genres mg ON m.id = mg.movie_id
		LEFT JOIN genres g ON mg.genre_id = g.id
		LEFT JOIN movie_artists ma ON m.id = ma.movie_id
		LEFT JOIN artists a ON ma.artist_id = a.id
		WHERE m.title LIKE ? OR m.description LIKE ? OR g.name LIKE ? OR a.name LIKE ?
		LIMIT ? OFFSET ?
	`

	// log.Printf("Executing query: %s\nWith parameters: %v, %v, %v, %v, %d, %d", queryString, query, query, query, query, limit, offset)

	rows, err := r.db.QueryContext(ctx, queryString, query, query, query, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.Duration,
			&movie.WatchURL,
			&movie.CreatedAt,
			&movie.UpdatedAt); err != nil {
			return nil, err
		}

		// Fetch genres for the movie
		genreRows, err := r.getGenresByMovieID(ctx, movie.ID)
		if err != nil {
			log.Fatalf("Error fetching genres: %v", err)
		}
		movie.Genres = genreRows
		movies = append(movies, movie)
	}

	return movies, nil
}

// GetGenresByMovieID retrieves genres associated with a given movie ID.
func (r *movieRepository) getGenresByMovieID(ctx context.Context, movieID string) ([]models.Genre, error) {
	query := `
		SELECT g.id, g.name
		FROM genres g
		JOIN movie_genres mg ON g.id = mg.genre_id
		WHERE mg.movie_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, movieID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(&genre.ID, &genre.Name); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		genres = append(genres, genre)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return genres, nil
}

func (r *movieRepository) TrackMovieView(ctx context.Context, movieID string) error {
	query := `
		INSERT INTO movie_views (movie_id, view_count, last_viewed_at)
		VALUES (?, 1, NOW())
		ON DUPLICATE KEY UPDATE 
			view_count = view_count + 1, 
			last_viewed_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, movieID)
	if err != nil {
		return err
	}
	return nil
}

// GetVoteByUserAndMovie checks if a user has already voted for a specific movie
func (r *movieRepository) GetVoteByUserAndMovie(ctx context.Context, userID, movieID string) (*models.Vote, error) {
	var vote models.Vote
	query := "SELECT id, user_id, movie_id FROM votes WHERE user_id = ? AND movie_id = ?"
	err := r.db.QueryRowContext(ctx, query, userID, movieID).Scan(&vote.ID, &vote.UserID, &vote.MovieID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &vote, nil
}

// CreateVote inserts a new vote into the database
func (r *movieRepository) CreateVote(ctx context.Context, userID, movieID string) error {
	query := "INSERT INTO votes (id, user_id, movie_id) VALUES (?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, uuid.NewString(), userID, movieID)
	return err
}
