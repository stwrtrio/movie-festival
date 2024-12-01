package repositories_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/stwrtrio/movie-festival/internal/repositories"
)

// Test for GetUserByUsername function in repository
// Integration test for GetUserByUsername function
func TestGetUserByUsername(t *testing.T) {
	// Create a user repository instance
	repo := repositories.NewUserRepository(testDB)

	// Test case: Create dummy user data
	userID := uuid.NewString()
	username := "testuserdummy"
	passwordHash := "$2a$10$7zOGb5S4F0TAMvuIEXJxH.yGjkoQ2I6ES4.l8P0e.mXJaX5aiRlYS" // Example bcrypt hash
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
