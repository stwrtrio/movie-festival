package repositories_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stwrtrio/movie-festival/internal/models"
	"github.com/stwrtrio/movie-festival/internal/repositories"
)

var testDB *sql.DB // Shared DB connection for all tests

func initTestDB() *sql.DB {
	// Load .env file in test setup

	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Return the DB connection
	return db
}

// TestMain is the entry point for tests
func TestMain(m *testing.M) {
	// Initialize test DB
	testDB = initTestDB()

	// Run tests
	code := m.Run()

	// Teardown: Close the DB connection
	testDB.Close()

	// Exit with test code
	os.Exit(code)
}

// TestCreateMovie_Success
func TestCreateMovie_Success(t *testing.T) {
	repo := repositories.NewMovieRepository(testDB)

	// Create a movie object
	movie := &models.Movie{
		Title:       "Inception",
		Description: "A mind-bending thriller",
		Duration:    148,
		Genres: []models.Genre{
			{Name: "Sci-fi"},
			{Name: "Action"},
		},
		WatchURL: "http://example.com/inception.mp4",
	}

	// insert the movie into the database
	err := repo.Create(context.Background(), movie)
	assert.NoError(t, err)

	// Verify the movie was inserted correctly
	var result models.Movie
	row := testDB.QueryRow("SELECT id, title FROM movies WHERE title = ?", movie.Title)
	err = row.Scan(&result.ID, &result.Title)

	assert.NoError(t, err)
	assert.Equal(t, movie.Title, result.Title)
}

// TestUpdateMovie_Success
func TestUpdateMovie_Success(t *testing.T) {
	repo := repositories.NewMovieRepository(testDB)

	// Create a movie object
	movie := &models.Movie{
		ID:          "4183d63a-e51a-4b24-9774-0c7051984071",
		Title:       "Inception 2",
		Description: "A mind-bending thriller",
		Duration:    148,
		Genres: []models.Genre{
			{Name: "Sci-fi"},
			{Name: "Action"},
		},
		WatchURL: "http://example.com/inception.mp4",
	}

	// update the movie into the database
	err := repo.Update(context.Background(), movie)

	assert.NoError(t, err)

	// Verify the movie was updated correctly
	var result models.Movie
	row := testDB.QueryRow("SELECT title FROM movies WHERE id = ?", movie.ID)
	err = row.Scan(&result.Title)

	assert.NoError(t, err)
	assert.Equal(t, movie.Title, result.Title)
}

func TestGetMostViewedMovie(t *testing.T) {
	// Mock database and sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Initialize repository
	repo := repositories.NewMovieRepository(db)

	// Test data
	movieID := uuid.NewString()
	movie := models.Movie{
		ID:          movieID,
		Title:       "Test Movie",
		Description: "A great movie",
		Duration:    120,
		WatchURL:    "http://example.com/movie.mp4",
		Views:       100,
		Genres: []models.Genre{
			{ID: 1, Name: "Action"},
			{ID: 2, Name: "Thriller"},
		},
	}

	// Mock the first query to get the most viewed movie
	movieQuery := `
		SELECT m.id, m.title, m.description, m.duration, m.watch_url, mv.view_count
		FROM movies m
		JOIN movie_views mv ON m.id = mv.movie_id
		ORDER BY mv.view_count DESC
		LIMIT 1
	`
	mock.ExpectQuery(movieQuery).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "duration", "watch_url", "view_count"}).
			AddRow(movie.ID, movie.Title, movie.Description, movie.Duration, movie.WatchURL, movie.Views))

	// Mock the second query to fetch genres for the movie
	genreQuery := `
		SELECT g.id, g.name
		FROM genres g
		JOIN movie_genres mg ON g.id = mg.genre_id
		WHERE mg.movie_id = ?
	`
	mock.ExpectQuery(genreQuery).
		WithArgs(movieID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "Action").
			AddRow(2, "Thriller"))

	// Call the repository function
	result, err := repo.GetMostViewedMovie(context.Background())

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, movie.ID, result.ID)
	assert.Equal(t, movie.Title, result.Title)
	assert.Equal(t, movie.Description, result.Description)
	assert.Equal(t, movie.Duration, result.Duration)
	assert.Equal(t, movie.WatchURL, result.WatchURL)
	assert.Equal(t, movie.Views, result.Views)
	assert.ElementsMatch(t, movie.Genres, result.Genres)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMostViewedGenre(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	// Create the movie repository with the mock database
	repo := repositories.NewMovieRepository(db)

	// Define expected results
	expectedGenreName := "Action"
	expectedGenreTotalViews := 500

	expectedQuery := `
		SELECT g.name, SUM(mv.view_count) AS total_views
		FROM movie_genres mg
		JOIN genres g ON mg.genre_id = g.id
		JOIN movie_views mv ON mg.movie_id = mv.movie_id
		GROUP BY g.id
		ORDER BY total_views DESC
		LIMIT 1
	`

	rows := sqlmock.NewRows([]string{"name", "total_views"}).
		AddRow(expectedGenreName, expectedGenreTotalViews)

	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnRows(rows)

	// Call the method under test
	genreName, totalViews, err := repo.GetMostViewedGenre(context.Background())

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, expectedGenreName, genreName)
	assert.Equal(t, expectedGenreTotalViews, totalViews)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
