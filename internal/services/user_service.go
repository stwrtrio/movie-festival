package services

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"github.com/stwrtrio/movie-festival/internal/helpers"
	"github.com/stwrtrio/movie-festival/internal/models"
	"github.com/stwrtrio/movie-festival/internal/repositories"
)

type UserService interface {
	Register(ctx context.Context, req models.RegisterRequest) error
	Login(ctx context.Context, username, password string) (string, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(ctx context.Context, req models.RegisterRequest) error {
	// Check if the username already exists
	existingUser, _ := s.repo.GetUserByUsername(ctx, req.Username)
	if existingUser != nil {
		return errors.New("username already exists")
	}

	// Create user model
	user := &models.User{
		ID:           uuid.NewString(),
		Username:     req.Username,
		PasswordHash: req.Password,
		Role:         "user",
	}

	// Save user in the repository
	return s.repo.CreateUser(ctx, user)
}

func (s *userService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	// Compare hashed passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := helpers.GenerateJWTToken(user.ID, user.Username, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
