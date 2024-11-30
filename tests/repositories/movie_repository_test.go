package repositories_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stwrtrio/movie-festival/internal/models"
	"github.com/stwrtrio/movie-festival/internal/repositories"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func cleanDummyData(movie *models.Movie) error {
	// Clean up
	_, err := testDB.Exec("DELETE FROM movies WHERE id = ?", movie.ID)
	if err != nil {
		return err
	}
	_, err = testDB.Exec("DELETE FROM movie_artists WHERE movie_id = ?", movie.ID)
	if err != nil {
		return err
	}
	_, err = testDB.Exec("DELETE FROM movie_genres WHERE movie_id = ?", movie.ID)
	if err != nil {
		return err
	}

	for _, artist := range movie.Artists {
		fmt.Println("artist_id:", artist.ID)
		_, err = testDB.Exec("DELETE FROM artists WHERE id = ?", artist.ID)
		if err != nil {
			return err
		}
	}

	for _, genre := range movie.Genres {
		fmt.Println("genre_id:", genre.ID)
		_, err = testDB.Exec("DELETE FROM genres WHERE id = ?", genre.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// TestCreateMovie_Success
func TestCreateMovie(t *testing.T) {
	// Create the repository
	repo := repositories.NewMovieRepository(testDB)

	// Create dummy data
	movie := &models.Movie{
		ID:          uuid.NewString(),
		Title:       "Test Movie",
		Description: "A test movie description",
		Duration:    120,
		WatchURL:    "http://example.com",
		Views:       0,
		Genres: []models.Genre{
			{Name: "Action"},
			{Name: "Adventure"},
		},
		Artists: []models.Artist{
			{ID: uuid.NewString(), Name: "John Doe"},
			{ID: uuid.NewString(), Name: "Jane Smith"},
		},
	}

	// Insert dummy data into the database
	err := repo.Create(context.Background(), movie)
	require.NoError(t, err)

	// Verify that the movie was created in the database
	var movieID string
	err = testDB.QueryRow("SELECT id FROM movies WHERE title = ?", movie.Title).Scan(&movieID)
	require.NoError(t, err)
	assert.NotEmpty(t, movieID)

	// Verify the genres were inserted and linked correctly
	for _, genre := range movie.Genres {
		var genreID int64
		err := testDB.QueryRow("SELECT id FROM genres WHERE name = ?", genre.Name).Scan(&genreID)
		require.NoError(t, err)
		assert.NotZero(t, genreID)

		// Verify the movie_genres link
		var movieGenreCount int
		err = testDB.QueryRow("SELECT COUNT(*) FROM movie_genres WHERE movie_id = ? AND genre_id = ?", movieID, genreID).Scan(&movieGenreCount)
		require.NoError(t, err)
		assert.Equal(t, 1, movieGenreCount)
	}

	// Verify the artists were inserted and linked correctly
	for _, artist := range movie.Artists {
		var artistID string
		err := testDB.QueryRow("SELECT id FROM artists WHERE name = ?", artist.Name).Scan(&artistID)
		require.NoError(t, err)
		assert.NotEmpty(t, artistID)

		// Verify the movie_artists link
		var movieArtistCount int
		err = testDB.QueryRow("SELECT COUNT(*) FROM movie_artists WHERE movie_id = ? AND artist_id = ?", movieID, artistID).Scan(&movieArtistCount)
		require.NoError(t, err)
		assert.Equal(t, 1, movieArtistCount)
	}

	// Clean up
	err = cleanDummyData(movie)
	require.NoError(t, err)
}

// TestUpdateMovie_Success
func TestUpdateMovie_Success(t *testing.T) {
	repo := repositories.NewMovieRepository(testDB)

	// Create dummy data to insert
	movie := &models.Movie{
		ID:          uuid.NewString(),
		Title:       "Test Movie",
		Description: "A test movie description",
		Duration:    120,
		WatchURL:    "http://example.com",
		Views:       0,
		Genres: []models.Genre{
			{Name: "Action"},
			{Name: "Adventure"},
		},
		Artists: []models.Artist{
			{ID: uuid.NewString(), Name: "John Doe"},
			{ID: uuid.NewString(), Name: "Jane Smith"},
		},
	}

	// Insert dummy data into the database
	err := repo.Create(context.Background(), movie)
	require.NoError(t, err)

	// Verify the movie was updated correctly
	var result models.Movie
	row := testDB.QueryRow("SELECT title FROM movies WHERE title = ?", movie.Title)
	err = row.Scan(&result.Title)
	assert.NoError(t, err)
	assert.Equal(t, movie.Title, result.Title)

	// Update dummy data
	movie.Title = "Test Movie Update"
	movie.Genres = []models.Genre{
		{Name: "Drama"},
		{Name: "Commedy"},
	}
	movie.Artists = []models.Artist{
		{ID: uuid.NewString(), Name: "John Smith"},
		{ID: uuid.NewString(), Name: "Jane Doe"},
	}

	// update the movie into the database
	err = repo.Update(context.Background(), movie)
	assert.NoError(t, err)

	// Verify the genres were inserted and linked correctly
	for _, genre := range movie.Genres {
		err := testDB.QueryRow("SELECT id FROM genres WHERE name = ?", genre.Name).Scan(&genre.ID)
		require.NoError(t, err)
		assert.NotZero(t, genre.ID)

		// Verify the movie_genres link
		var movieGenreCount int
		err = testDB.QueryRow("SELECT COUNT(*) FROM movie_genres WHERE movie_id = ? AND genre_id = ?", movie.ID, genre.ID).Scan(&movieGenreCount)
		require.NoError(t, err)
		assert.Equal(t, 1, movieGenreCount)
	}

	// Verify the artists were inserted and linked correctly
	for _, artist := range movie.Artists {
		var artistID string
		err := testDB.QueryRow("SELECT id FROM artists WHERE name = ?", artist.Name).Scan(&artistID)
		require.NoError(t, err)
		assert.NotEmpty(t, artistID)

		// Verify the movie_artists link
		var movieArtistCount int
		err = testDB.QueryRow("SELECT COUNT(*) FROM movie_artists WHERE movie_id = ? AND artist_id = ?", movie.ID, artistID).Scan(&movieArtistCount)
		require.NoError(t, err)
		assert.Equal(t, 1, movieArtistCount)
	}

	// Clean up
	err = cleanDummyData(movie)
	require.NoError(t, err)
}

func TestGetMostViewedMovie(t *testing.T) {
	// Create a new repository instance with the test database
	repo := repositories.NewMovieRepository(testDB)

	// Test data
	movieID := uuid.NewString()
	movie := &models.Movie{
		ID:          movieID,
		Title:       "Test Movie",
		Description: "A great movie",
		Duration:    120,
		WatchURL:    "http://example.com/movie.mp4",
		Genres: []models.Genre{
			{ID: 1, Name: "Action"},
			{ID: 2, Name: "Thriller"},
		},
	}

	// Insert dummy data into the database
	err := repo.Create(context.Background(), movie)
	assert.NoError(t, err)

	// Verify the movie was inserted
	var res models.Movie
	row := testDB.QueryRow("SELECT id FROM movies WHERE id = ?", movie.ID)
	err = row.Scan(&res.ID)
	assert.NoError(t, err)
	assert.Equal(t, movie.ID, res.ID)

	// Set up some dummy data for the test
	_, err = testDB.Exec(`
		INSERT INTO movie_views (movie_id, view_count, last_viewed_at)
		VALUES (?, 5000, NOW())
	`, movieID)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
		assert.NoError(t, err)
	}

	// Call the repository function
	result, err := repo.GetMostViewedMovie(context.Background())

	// Assert the result
	assert.NoError(t, err)
	assert.NotNil(t, movie)
	assert.Equal(t, result.Title, movie.Title)

	// Clean up test data from the database
	err = cleanDummyData(movie)
	require.NoError(t, err)

	_, err = testDB.Exec("DELETE FROM movie_views WHERE movie_id = ?", movieID)
	require.NoError(t, err)
}
