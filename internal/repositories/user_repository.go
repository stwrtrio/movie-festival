package repositories

import (
	"context"
	"database/sql"

	"github.com/stwrtrio/movie-festival/internal/models"
)

type UserRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, password_hash, role FROM users WHERE username = ?"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
