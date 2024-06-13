package util

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// AccessTokenClaims for access tokens
type AccessTokenClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// RefreshTokenClaims for refresh tokens
type RefreshTokenClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateAccessToken generates a new JWT access token
func GenerateAccessToken(userID uint, jwtSecret string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute) // 15 minutes expiration
	claims := &AccessTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			// You can add other registered claims if needed
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// GenerateRefreshToken generates a new JWT refresh token
func GenerateRefreshToken(userID uint, jwtSecret string) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour) // 7 days expiration
	claims := &RefreshTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			// You can add other registered claims if needed
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return refreshToken.SignedString([]byte(jwtSecret))
}

// ValidateToken validates a JWT token (both access and refresh)
func ValidateToken(tokenString string, jwtSecret string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
}
