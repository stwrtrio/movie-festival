package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/stwrtrio/movie-festival/config"
	"github.com/stwrtrio/movie-festival/internal/helpers"
	"github.com/stwrtrio/movie-festival/internal/utils"
)

func GetUserFromContext(c echo.Context) (*helpers.Claims, bool) {
	claims, ok := c.Get("user").(*helpers.Claims)
	return claims, ok
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := validateToken(c)
		if err != nil {
			return utils.FailResponse(c, http.StatusUnauthorized, err.Error())
		}

		// Store user claims in Echo's context
		c.Set("user", claims)

		// Proceed to the next handler
		return next(c)
	}
}

// AdminAuthMiddleware is a middleware to check if the user has an admin role.
func AdminAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := validateToken(c)
		if err != nil {
			return utils.FailResponse(c, http.StatusUnauthorized, err.Error())
		}

		// Check if the user role is admin
		if claims.Role != "admin" {
			return utils.FailResponse(c, http.StatusForbidden, "Access denied. Admins only.")
		}

		// Store user claims in Echo's context (optional for admin routes)
		c.Set("user", claims)

		// Proceed to the next handler
		return next(c)
	}
}

func validateToken(c echo.Context) (*helpers.Claims, error) {
	// Get the Authorization header
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header is missing")
	}

	// Check if the token has "Bearer" prefix
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return nil, errors.New("invalid Authorization format")
	}

	// Extract the token
	token := tokenParts[1]

	// Validate the token
	claims, err := helpers.ValidateJWTToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	// Check if the token is blacklisted in Redis
	redisClient := config.RedisClient
	ctx := c.Request().Context()
	isBlacklisted, err := redisClient.Get(ctx, claims.JTI).Result()
	if err == nil && isBlacklisted == "true" {
		return nil, errors.New("token has been revoked")
	}

	return claims, nil
}
