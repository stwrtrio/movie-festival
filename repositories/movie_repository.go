package repositories

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/stwrtrio/movie-festival/models"
)

type MovieRepository interface {
	Create(movie *models.Movie) error
}

type movieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) MovieRepository {
	return &movieRepository{db}
}

func (r *movieRepository) Create(movie *models.Movie) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// create new id
	movie.ID = uuid.NewString()

	// Insert movie
	query := `
        INSERT INTO movies (id, title, description, duration, genres, watch_url) 
        VALUES (?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(query, movie.ID, movie.Title, movie.Description, movie.Duration, movie.Genres, movie.WatchURL)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert artists and associate with the movie
	for _, artist := range movie.Artists {
		var artistID string
		err := tx.QueryRow("SELECT id FROM artists WHERE name = ?", artist.Name).Scan(&artistID)
		if err == sql.ErrNoRows {
			// Artist doesn't exist
			artistID = uuid.NewString()
			_, err = tx.Exec("INSERT INTO artists (id, name) VALUES (?, ?)", artistID, artist.Name)
			if err != nil {
				tx.Rollback()
				return err
			}

		} else if err != nil {
			tx.Rollback()
			return err
		}

		// Associate artist with the movie
		_, err = tx.Exec("INSERT INTO movie_artists (movie_id, artist_id) VALUES (?, ?)", movie.ID, artistID)
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
