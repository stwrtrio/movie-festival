package middlewares

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/stwrtrio/movie-festival/config"
	"github.com/stwrtrio/movie-festival/internal/helpers"
	"github.com/stwrtrio/movie-festival/internal/utils"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return utils.FailResponse(c, http.StatusUnauthorized, "Authorization header is missing")
		}

		// Check if the token has "Bearer" prefix
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return utils.FailResponse(c, http.StatusUnauthorized, "Invalid Authorization format")
		}

		// Extract the token
		token := tokenParts[1]

		// Validate the token
		claims, err := helpers.ValidateJWTToken(token)
		if err != nil {
			return utils.FailResponse(c, http.StatusUnauthorized, "Invalid or expired token")
		}

		// Check if the token is blacklisted in Redis
		redisClient := config.RedisClient
		ctx := c.Request().Context()

		isBlacklisted, err := redisClient.Get(ctx, claims.JTI).Result()
		if err == nil && isBlacklisted == "true" {
			return utils.FailResponse(c, http.StatusUnauthorized, "Token has been revoked")
		}

		// Store user claims in Echo's context
		c.Set("user", claims)

		// Proceed to the next handler
		return next(c)
	}
}

func GetUserFromContext(c echo.Context) (*helpers.Claims, bool) {
	claims, ok := c.Get("user").(*helpers.Claims)
	return claims, ok
}
