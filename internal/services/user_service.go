package services

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/stwrtrio/movie-festival/internal/helpers"
	"github.com/stwrtrio/movie-festival/internal/repositories"
)

type UserService interface {
	Login(ctx context.Context, username, password string) (string, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}

	fmt.Println("here---")
	fmt.Println("user hash: ", user.PasswordHash)
	fmt.Println("user pass: ", password)

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
