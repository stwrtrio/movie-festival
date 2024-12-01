package repositories_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/stwrtrio/movie-festival/internal/models"
	"github.com/stwrtrio/movie-festival/internal/repositories"
)

// Test for CreateUser function in repository
func TestCreateUser(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name          string
		user          *models.User
		expectedError string
	}{
		{
			name: "Successful user creation",
			user: &models.User{
				ID:           uuid.NewString(),
				Username:     "testuser",
				PasswordHash: "password123",
				Role:         "user",
			},
			expectedError: "",
		},
		{
			name: "Duplicate username",
			user: &models.User{
				ID:           uuid.NewString(),
				Username:     "duplicateuser",
				PasswordHash: "password123",
				Role:         "user",
			},
			expectedError: "Error 1062",
		},
	}

	// Create a repository instance
	userRepo := repositories.NewUserRepository(testDB)

	// Insert a dummy user to test duplicate handling
	dummyUsername := "duplicateuser"

	// Cleanup the dummy user if exist
	_, err := testDB.Exec("DELETE FROM users WHERE username = ?", dummyUsername)
	assert.NoError(t, err)

	// Generate password hash and then insert dummy user data for duplicate condition
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	_, err = testDB.Exec("INSERT INTO users (id, username, password_hash, role) VALUES (?, ?, ?, ?)",
		uuid.NewString(), dummyUsername, hashedPassword, "user")
	assert.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Hash the password for the test case
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tc.user.PasswordHash), bcrypt.DefaultCost)
			assert.NoError(t, err)
			tc.user.PasswordHash = string(hashedPassword)

			// Execute the CreateUser function
			err = userRepo.CreateUser(context.Background(), tc.user)

			if tc.expectedError != "" {
				// Validate that the error message matches expectations
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)

				// Verify the user was successfully inserted
				var count int
				query := "SELECT COUNT(*) FROM users WHERE username = ?"
				err := testDB.QueryRow(query, tc.user.Username).Scan(&count)
				assert.NoError(t, err)
				assert.Equal(t, 1, count)

				// Cleanup: Remove the inserted user
				_, err = testDB.Exec("DELETE FROM users WHERE username = ?", tc.user.Username)
				assert.NoError(t, err)
			}
		})
	}

	// Cleanup the dummy user
	_, err = testDB.Exec("DELETE FROM users WHERE username = ?", dummyUsername)
	assert.NoError(t, err)
}

// Test for GetUserByUsername function in repository
// Integration test for GetUserByUsername function
func TestGetUserByUsername(t *testing.T) {
	// Create a user repository instance
	repo := repositories.NewUserRepository(testDB)

	// Test case: Create dummy user data
	userID := uuid.NewString()
	username := "testuserdummy"
	passwordHash := "$2a$10$7zOGb5S4F0TAMvuIEXJxH.yGjkoQ2I6ES4.l8P0e.mXJaX5aiRlYS" // Example bcrypt hash for password123
	role := "user"

	// Insert dummy data into the 'users' table
	insertQuery := "INSERT INTO users (id,username, password_hash, role) VALUES (?, ?, ?, ?)"
	_, err := testDB.Exec(insertQuery, userID, username, passwordHash, role)
	assert.NoError(t, err, "Error inserting dummy user")

	// Verify: Call GetUserByUsername and check the result
	t.Run("Verify user is fetched", func(t *testing.T) {
		user, err := repo.GetUserByUsername(context.Background(), username)
		assert.NoError(t, err, "Error fetching user by username")
		assert.NotNil(t, user, "User should not be nil")
		assert.Equal(t, username, user.Username, "Usernames should match")
		assert.Equal(t, role, user.Role, "Roles should match")
	})

	// Clean Up: Delete the inserted dummy user from the database
	t.Run("Clean up inserted data", func(t *testing.T) {
		deleteQuery := "DELETE FROM users WHERE id = ?"
		_, err := testDB.Exec(deleteQuery, userID)
		assert.NoError(t, err, "Error cleaning up inserted data")
	})

}
