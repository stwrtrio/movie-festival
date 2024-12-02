package helpers

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	JTI      string `json:"jti"`
	jwt.RegisteredClaims
}

// LoadJWTExpiry load JWT_EXPIRY in .env
func LoadJWTExpiry() time.Duration {
	expiryStr := os.Getenv("JWT_EXPIRY")
	expiry, err := time.ParseDuration(expiryStr)
	if err != nil {
		log.Fatalf("Invalid JWT_EXPIRY in .env: %v", err)
	}
	return expiry
}

// GenerateJWTToken creates a JWT token for the user
func GenerateJWTToken(userID string, username, role string) (string, error) {
	expiry := LoadJWTExpiry()
	expirationTime := time.Now().Add(expiry)
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		JTI: uuid.NewString(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secretKey))
}

// ValidateJWTToken parses and validates a JWT token
func ValidateJWTToken(tokenString string) (*Claims, error) {
	secretKey := os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}
