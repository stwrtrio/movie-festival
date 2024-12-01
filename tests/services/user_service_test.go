package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/stwrtrio/movie-festival/internal/models"
	"github.com/stwrtrio/movie-festival/internal/services"
	"github.com/stwrtrio/movie-festival/tests/mocks"
)

func TestRegister(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name          string
		request       models.RegisterRequest
		mockSetup     func(mockRepo *mocks.MockUserRepository)
		expectedError string
	}{
		{
			name: "Successful registration",
			request: models.RegisterRequest{
				Username: "newuser",
				Password: "securepassword",
			},
			mockSetup: func(mockRepo *mocks.MockUserRepository) {
				// Mock GetUserByUsername to return nil, indicating the username does not exist
				mockRepo.EXPECT().GetUserByUsername(gomock.Any(), "newuser").Return(nil, nil)

				// Mock CreateUser to succeed
				mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "Username already exists",
			request: models.RegisterRequest{
				Username: "existinguser",
				Password: "securepassword",
			},
			mockSetup: func(mockRepo *mocks.MockUserRepository) {
				// Mock GetUserByUsername to return an existing user
				mockRepo.EXPECT().GetUserByUsername(gomock.Any(), "existinguser").Return(&models.User{
					ID:           uuid.NewString(),
					Username:     "existinguser",
					PasswordHash: "hashedpassword",
					Role:         "user",
				}, nil)
			},
			expectedError: "username already exists",
		},
		{
			name: "Repository error on CreateUser",
			request: models.RegisterRequest{
				Username: "newuser",
				Password: "securepassword",
			},
			mockSetup: func(mockRepo *mocks.MockUserRepository) {
				// Mock GetUserByUsername to return nil
				mockRepo.EXPECT().GetUserByUsername(gomock.Any(), "newuser").Return(nil, nil)

				// Mock CreateUser to return an error
				mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(errors.New("repository error"))
			},
			expectedError: "repository error",
		},
	}

	// Iterate through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create a mock repository
			mockRepo := mocks.NewMockUserRepository(ctrl)

			// Set up mock behavior
			tc.mockSetup(mockRepo)

			// Create the service
			userService := services.NewUserService(mockRepo)

			// Call the service method
			err := userService.Register(context.Background(), tc.request)

			// Validate results
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserLogin(t *testing.T) {
	// Test cases
	testCases := []struct {
		name          string
		username      string
		password      string
		mockUser      *models.User
		mockSetup     func(mockRepo *mocks.MockUserRepository)
		mockError     error
		expectedToken string
		expectedError string
	}{
		{
			name:     "Valid credentials",
			username: "user123",
			password: "password123",
			mockUser: &models.User{
				ID:           uuid.NewString(),
				Username:     "user123",
				PasswordHash: "$2a$10$7zOGb5S4F0TAMvuIEXJxH.yGjkoQ2I6ES4.l8P0e.mXJaX5aiRlYS", // bcrypt hash for "password123"
				Role:         "user",
			},
			mockSetup: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().
					GetUserByUsername(gomock.Any(), "user123").
					Return(&models.User{
						ID:           uuid.NewString(),
						Username:     "user123",
						PasswordHash: "$2a$10$7zOGb5S4F0TAMvuIEXJxH.yGjkoQ2I6ES4.l8P0e.mXJaX5aiRlYS", // bcrypt hash for "password123"
						Role:         "user",
					}, nil)
			},
			expectedToken: "non-empty", // Check for a non-empty token
			expectedError: "",
		},
		{
			name:     "Invalid password",
			username: "user123",
			password: "wrongpassword",
			mockUser: &models.User{
				ID:           uuid.NewString(),
				Username:     "user123",
				PasswordHash: "$2a$10$7zOGb5S4F0TAMvuIEXJxH.yGjkoQ2I6ES4.l8P0e.mXJaX5aiRlYS",
				Role:         "user",
			},
			mockSetup: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().
					GetUserByUsername(gomock.Any(), "user123").
					Return(&models.User{
						ID:           uuid.NewString(),
						Username:     "user123",
						PasswordHash: "$2a$10$7zOGb5S4F0TAMvuIEXJxH.yGjkoQ2I6ES4.l8P0e.mXJaX5aiRlYS",
						Role:         "user",
					}, nil)
			},
			expectedToken: "",
			expectedError: "invalid credentials",
		},
		{
			name:     "User not found",
			username: "user123",
			password: "password123",
			mockUser: nil,
			mockSetup: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().
					GetUserByUsername(gomock.Any(), "user123").
					Return(nil, nil)
			},
			expectedToken: "",
			expectedError: "invalid credentials",
		},
		{
			name:     "Repository error",
			username: "user123",
			password: "password123",
			mockUser: nil,
			mockSetup: func(mockRepo *mocks.MockUserRepository) {
				mockRepo.EXPECT().
					GetUserByUsername(gomock.Any(), "user123").
					Return(nil, errors.New("database error"))
			},
			expectedToken: "",
			expectedError: "database error",
		},
	}

	// Iterate through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new gomock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create a mock repository
			mockRepo := mocks.NewMockUserRepository(ctrl)

			// Setup the mock behavior for the current test case
			tc.mockSetup(mockRepo)

			// Create the service
			userService := services.NewUserService(mockRepo)

			// Call the service method
			token, err := userService.Login(context.Background(), tc.username, tc.password)

			// Validate results
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
			}

			if tc.expectedToken == "non-empty" {
				assert.NotEmpty(t, token)
			} else {
				assert.Equal(t, tc.expectedToken, token)
			}
		})
	}
}
