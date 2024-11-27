package repositories_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stwrtrio/movie-festival/models"
	"github.com/stwrtrio/movie-festival/repositories"
)

var testDB *sql.DB // Shared DB connection for all tests

func initTestDB() *sql.DB {
	// Load .env file in test setup

	if err := godotenv.Load("../.env"); err != nil {
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
		Genres:      "Sci-Fi",
		WatchURL:    "http://example.com/inception.mp4",
	}

	// insert the movie into the database
	err := repo.Create(movie)

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
		Genres:      "Sci-Fi",
		WatchURL:    "http://example.com/inception.mp4",
	}

	// update the movie into the database
	err := repo.Update(movie)

	assert.NoError(t, err)

	// Verify the movie was updated correctly
	var result models.Movie
	row := testDB.QueryRow("SELECT title FROM movies WHERE id = ?", movie.ID)
	err = row.Scan(&result.Title)

	assert.NoError(t, err)
	assert.Equal(t, movie.Title, result.Title)
}
