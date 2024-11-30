package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/stwrtrio/movie-festival/internal/models"
	"github.com/stwrtrio/movie-festival/internal/services"
	"github.com/stwrtrio/movie-festival/tests/mocks"
)

func TestCreateMovie(t *testing.T) {

	// Test cases
	tests := []struct {
		name          string
		inputMovie    *models.Movie
		mockSetup     func(mockRepo *mocks.MockMovieRepository)
		expectedError error
	}{
		{
			name: "Success - Movie created successfully",
			inputMovie: &models.Movie{
				ID:          uuid.NewString(),
				Title:       "Test Movie",
				Description: "A test movie description",
				Duration:    120,
				WatchURL:    "http://test.com",
				Genres: []models.Genre{
					{Name: "Action"},
					{Name: "Drama"},
				},
				Artists: []models.Artist{
					{ID: uuid.NewString(), Name: "John Doe"},
					{ID: uuid.NewString(), Name: "Jane Doe"},
				},
			},
			mockSetup: func(mockRepo *mocks.MockMovieRepository) {
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Failure - Repository returns an error",
			inputMovie: &models.Movie{
				Title:       "Test Movie",
				Description: "A test movie description",
				Duration:    120,
				WatchURL:    "http://test.com",
				Artists: []models.Artist{
					{ID: uuid.NewString(), Name: "John Doe"},
					{ID: uuid.NewString(), Name: "Jane Doe"},
				},
			},
			mockSetup: func(mockRepo *mocks.MockMovieRepository) {
				mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("repository error"))
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new gomock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create a mock repository
			mockRepo := mocks.NewMockMovieRepository(ctrl)
			tt.mockSetup(mockRepo)

			// Mock Redis client
			mockRedisClient, _ := redismock.NewClientMock()

			// Create the service
			movieService := services.NewMovieService(mockRepo, mockRedisClient)

			// Execute the service method
			err := movieService.CreateMovie(context.TODO(), tt.inputMovie)

			// Assert the result
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateMovie(t *testing.T) {
	// Test cases
	tests := []struct {
		name          string
		inputMovie    *models.Movie
		mockSetup     func(mockRepo *mocks.MockMovieRepository)
		expectedError error
	}{
		{
			name: "Success - Movie updated successfully",
			inputMovie: &models.Movie{
				ID:          "movie1",
				Title:       "Updated Movie",
				Description: "Updated description",
				Duration:    150,
				WatchURL:    "http://updated.com",
			},
			mockSetup: func(mockRepo *mocks.MockMovieRepository) {
				mockRepo.EXPECT().FindMovieByID(gomock.Any(), "movie1").Return(models.Movie{ID: "movie1"}, nil)
				// Mock Update to return no error
				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Failure - Movie does not exist",
			inputMovie: &models.Movie{
				ID:          "nonexistent-id",
				Title:       "Nonexistent Movie",
				Description: "This movie does not exist",
				Duration:    120,
				WatchURL:    "http://nonexistent-movie.com",
			},
			mockSetup: func(mockRepo *mocks.MockMovieRepository) {
				mockRepo.EXPECT().
					FindMovieByID(gomock.Any(), "nonexistent-id").
					Return(models.Movie{}, nil)
			},
			expectedError: errors.New("service err: movie is not exists"),
		},
		{
			name: "Failure - Error finding movie",
			inputMovie: &models.Movie{
				ID: "movie3",
			},
			mockSetup: func(mockRepo *mocks.MockMovieRepository) {
				mockRepo.EXPECT().FindMovieByID(gomock.Any(), gomock.Any()).Return(models.Movie{}, errors.New("repository error"))
			},
			expectedError: errors.New("repository error"),
		},
		{
			name: "Failure - Error updating movie",
			inputMovie: &models.Movie{
				ID: "movie4",
			},
			mockSetup: func(mockRepo *mocks.MockMovieRepository) {
				mockRepo.EXPECT().FindMovieByID(gomock.Any(), "movie4").Return(models.Movie{ID: "movie4"}, nil)
				mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("update error"))
			},
			expectedError: errors.New("update error"),
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new gomock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create a mock repository
			mockRepo := mocks.NewMockMovieRepository(ctrl)
			tt.mockSetup(mockRepo)

			// Mock Redis client
			mockRedisClient, _ := redismock.NewClientMock()

			// Initialize the movie service with the mock repository
			movieService := services.NewMovieService(mockRepo, mockRedisClient)

			// Call the UpdateMovie service method
			err := movieService.UpdateMovie(context.TODO(), tt.inputMovie)

			// Check if the error matches the expected error
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetMostViewedMovie(t *testing.T) {
	// Define test cases
	tests := []struct {
		name          string
		mockSetup     func(mockRepo *mocks.MockMovieRepository)
		expectedMovie *models.Movie
		expectedError error
	}{
		{
			name: "Success - Most viewed movie retrieved",
			mockSetup: func(mockRepo *mocks.MockMovieRepository) {
				mockRepo.EXPECT().
					GetMostViewedMovie(gomock.Any()).
					Return(&models.Movie{
						ID:          "movie1",
						Title:       "Most Viewed Movie",
						Description: "A popular movie",
						Duration:    120,
						WatchURL:    "http://most-viewed.com",
						Views:       100,
						Genres: []models.Genre{
							{Name: "Action"},
							{Name: "Adventure"},
						},
					}, nil)
			},
			expectedMovie: &models.Movie{
				ID:          "movie1",
				Title:       "Most Viewed Movie",
				Description: "A popular movie",
				Duration:    120,
				WatchURL:    "http://most-viewed.com",
				Views:       100,
				Genres: []models.Genre{
					{Name: "Action"},
					{Name: "Adventure"},
				},
			},
			expectedError: nil,
		},
		{
			name: "Failure - Repository returns an error",
			mockSetup: func(mockRepo *mocks.MockMovieRepository) {
				mockRepo.EXPECT().
					GetMostViewedMovie(gomock.Any()).
					Return(nil, errors.New("repository error"))
			},
			expectedMovie: nil,
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new gomock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create a mock repository
			mockRepo := mocks.NewMockMovieRepository(ctrl)
			tt.mockSetup(mockRepo)

			// Create the service
			movieService := services.NewMovieService(mockRepo, nil)

			// Execute the service method
			movie, err := movieService.GetMostViewedMovie(context.TODO())

			// Assert the result
			if tt.expectedError != nil {
				assert.Nil(t, movie)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NotNil(t, movie)
				assert.Equal(t, tt.expectedMovie, movie)
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetMostViewedGenre(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		page           int
		pageSize       int
		sortOrder      string
		mockSetup      func(mockRepo *mocks.MockMovieRepository)
		expectedResult []models.GenreView
		expectedError  error
	}{
		{
			name:      "Success - Most viewed genres retrieved (DESC order)",
			page:      1,
			pageSize:  5,
			sortOrder: "DESC",
			mockSetup: func(mockRepo *mocks.MockMovieRepository) {
				mockRepo.EXPECT().
					GetMostViewedGenre(gomock.Any(), 1, 5, "DESC").
					Return([]models.GenreView{
						{Name: "Action", ViewCount: 150},
						{Name: "Adventure", ViewCount: 120},
					}, nil)
			},
			expectedResult: []models.GenreView{
				{Name: "Action", ViewCount: 150},
				{Name: "Adventure", ViewCount: 120},
			},
			expectedError: nil,
		},
		{
			name:      "Success - Most viewed genres retrieved (ASC order)",
			page:      2,
			pageSize:  3,
			sortOrder: "ASC",
			mockSetup: func(mockRepo *mocks.MockMovieRepository) {
				mockRepo.EXPECT().
					GetMostViewedGenre(gomock.Any(), 2, 3, "ASC").
					Return([]models.GenreView{
						{Name: "Horror", ViewCount: 50},
						{Name: "Comedy", ViewCount: 30},
					}, nil)
			},
			expectedResult: []models.GenreView{
				{Name: "Horror", ViewCount: 50},
				{Name: "Comedy", ViewCount: 30},
			},
			expectedError: nil,
		},
		{
			name:      "Failure - Invalid sortOrder defaults to DESC",
			page:      1,
			pageSize:  5,
			sortOrder: "INVALID",
			mockSetup: func(mockRepo *mocks.MockMovieRepository) {
				mockRepo.EXPECT().
					GetMostViewedGenre(gomock.Any(), 1, 5, "DESC").
					Return([]models.GenreView{
						{Name: "Sci-Fi", ViewCount: 100},
					}, nil)
			},
			expectedResult: []models.GenreView{
				{Name: "Sci-Fi", ViewCount: 100},
			},
			expectedError: nil,
		},
		{
			name:      "Failure - Repository returns error",
			page:      1,
			pageSize:  5,
			sortOrder: "DESC",
			mockSetup: func(mockRepo *mocks.MockMovieRepository) {
				mockRepo.EXPECT().
					GetMostViewedGenre(gomock.Any(), 1, 5, "DESC").
					Return(nil, errors.New("repository error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new gomock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create a mock repository
			mockRepo := mocks.NewMockMovieRepository(ctrl)
			tt.mockSetup(mockRepo)

			// Create the service
			movieService := services.NewMovieService(mockRepo, nil)

			// Execute the service method
			result, err := movieService.GetMostViewedGenre(context.TODO(), tt.page, tt.pageSize, tt.sortOrder)

			// Assert the result
			if tt.expectedError != nil {
				assert.Nil(t, result)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult, result)
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetAllMovies(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		limit          int
		offset         int
		mockSetup      func(mockRepo *mocks.MockMovieRepository)
		expectedResult []models.Movie
		expectedError  error
	}{
		{
			name:   "Success - Movies retrieved successfully",
			limit:  5,
			offset: 0,
			mockSetup: func(mockRepo *mocks.MockMovieRepository) {
				mockRepo.EXPECT().
					GetAllMovies(gomock.Any(), 5, 0).
					Return([]models.Movie{
						{
							ID:          "movie1",
							Title:       "Movie One",
							Description: "A test movie one",
							Duration:    120,
							WatchURL:    "http://movie1.com",
						},
						{
							ID:          "movie2",
							Title:       "Movie Two",
							Description: "A test movie two",
							Duration:    90,
							WatchURL:    "http://movie2.com",
						},
					}, nil)
			},
			expectedResult: []models.Movie{
				{
					ID:          "movie1",
					Title:       "Movie One",
					Description: "A test movie one",
					Duration:    120,
					WatchURL:    "http://movie1.com",
				},
				{
					ID:          "movie2",
					Title:       "Movie Two",
					Description: "A test movie two",
					Duration:    90,
					WatchURL:    "http://movie2.com",
				},
			},
			expectedError: nil,
		},
		{
			name:   "Success - No movies available",
			limit:  5,
			offset: 0,
			mockSetup: func(mockRepo *mocks.MockMovieRepository) {
				mockRepo.EXPECT().
					GetAllMovies(gomock.Any(), 5, 0).
					Return([]models.Movie{}, nil)
			},
			expectedResult: []models.Movie{},
			expectedError:  nil,
		},
		{
			name:   "Failure - Repository error",
			limit:  10,
			offset: 5,
			mockSetup: func(mockRepo *mocks.MockMovieRepository) {
				mockRepo.EXPECT().
					GetAllMovies(gomock.Any(), 10, 5).
					Return(nil, errors.New("repository error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new gomock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create a mock repository
			mockRepo := mocks.NewMockMovieRepository(ctrl)
			tt.mockSetup(mockRepo)

			// Create the service
			movieService := services.NewMovieService(mockRepo, nil)

			// Execute the service method
			result, err := movieService.GetAllMovies(context.TODO(), tt.limit, tt.offset)

			// Assert the result
			if tt.expectedError != nil {
				assert.Nil(t, result)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult, result)
				assert.NoError(t, err)
			}
		})
	}
}

func TestSearchMoviesService(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		query          string
		limit          int
		offset         int
		mockRepoSetup  func(mockRepo *mocks.MockMovieRepository)
		expectedResult []models.Movie
		expectedError  error
	}{
		{
			name:   "Success - Movies found",
			query:  "action",
			limit:  5,
			offset: 0,
			mockRepoSetup: func(mockRepo *mocks.MockMovieRepository) {
				mockRepo.EXPECT().
					SearchMovies(gomock.Any(), "action", 5, 0).
					Return([]models.Movie{
						{
							ID:          "movie1",
							Title:       "Action Movie 1",
							Description: "A great action movie",
							Duration:    120,
							WatchURL:    "http://actionmovie1.com",
						},
						{
							ID:          "movie2",
							Title:       "Action Movie 2",
							Description: "Another action-packed movie",
							Duration:    130,
							WatchURL:    "http://actionmovie2.com",
						},
					}, nil)
			},
			expectedResult: []models.Movie{
				{
					ID:          "movie1",
					Title:       "Action Movie 1",
					Description: "A great action movie",
					Duration:    120,
					WatchURL:    "http://actionmovie1.com",
				},
				{
					ID:          "movie2",
					Title:       "Action Movie 2",
					Description: "Another action-packed movie",
					Duration:    130,
					WatchURL:    "http://actionmovie2.com",
				},
			},
			expectedError: nil,
		},
		{
			name:   "Error - No movies found",
			query:  "nonexistent",
			limit:  5,
			offset: 0,
			mockRepoSetup: func(mockRepo *mocks.MockMovieRepository) {
				mockRepo.EXPECT().
					SearchMovies(gomock.Any(), "nonexistent", 5, 0).
					Return([]models.Movie{}, nil) // Return an empty slice
			},
			expectedResult: []models.Movie{}, // Expect an empty slice
			expectedError:  nil,
		},
		{
			name:   "Error - Repository returns error",
			query:  "action",
			limit:  5,
			offset: 0,
			mockRepoSetup: func(mockRepo *mocks.MockMovieRepository) {
				mockRepo.EXPECT().
					SearchMovies(gomock.Any(), "action", 5, 0).
					Return(nil, errors.New("repository error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new gomock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create a mock repository
			mockRepo := mocks.NewMockMovieRepository(ctrl)
			tt.mockRepoSetup(mockRepo)

			// Create the service
			movieService := services.NewMovieService(mockRepo, nil) // Assuming no Redis for now

			// Execute the service method
			result, err := movieService.SearchMovies(context.TODO(), tt.query, tt.limit, tt.offset)

			// Assert the results
			if tt.expectedError != nil {
				assert.Nil(t, result)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult, result)
				assert.NoError(t, err)
			}
		})
	}
}

func TestTrackMovieViewService(t *testing.T) {
	// Define test cases
	tests := []struct {
		name          string
		movieID       string
		mockRepoSetup func(mockRepo *mocks.MockMovieRepository)
		expectedError error
	}{
		{
			name:    "Success - Movie view tracked",
			movieID: "movie123",
			mockRepoSetup: func(mockRepo *mocks.MockMovieRepository) {
				mockRepo.EXPECT().
					TrackMovieView(gomock.Any(), "movie123").
					Return(nil) // No error on success
			},
			expectedError: nil,
		},
		{
			name:    "Error - Movie view tracking fails",
			movieID: "movie456",
			mockRepoSetup: func(mockRepo *mocks.MockMovieRepository) {
				mockRepo.EXPECT().
					TrackMovieView(gomock.Any(), "movie456").
					Return(errors.New("repository error")) // Return error from repository
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new gomock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create a mock repository
			mockRepo := mocks.NewMockMovieRepository(ctrl)
			tt.mockRepoSetup(mockRepo)

			// Create the service
			movieService := services.NewMovieService(mockRepo, nil) // Assuming no Redis for now

			// Execute the service method
			err := movieService.TrackMovieView(context.TODO(), tt.movieID)

			// Assert the results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
