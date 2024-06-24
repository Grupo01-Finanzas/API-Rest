package middleware

import (
	"ApiRestFinance/internal/model/entities/enums"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// AuthMiddleware is a JWT authentication middleware for Gin
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		// Check if the token is in the format "Bearer {token}"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			return
		}

		tokenString := tokenParts[1]

		// Parse and validate the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Extract claims and set them in the context
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Unable to extract claims"})
			return
		}
		c.Set("claims", claims)

		// Extract user ID from claims
		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Unable to extract user ID"})
			return
		}
		userIDUint := uint(userID)
		c.Set("user_id", userIDUint)
		// Extract user role from claims
		rol := claims["rol"].(string)

		if rol == "" {
			fmt.Println("Rol is empty")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Unable to extract user role"})
			return
		}

		role := enums.Role(rol)

		c.Set("rol", role)
		c.Next()

	}
}

func GetUserIDFromContext(ctx *gin.Context) uint {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return 0
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		return 0
	}
	return userIDUint
}

func GetUserRoleFromContext(c *gin.Context) enums.Role {
	value, exists := c.Get("rol")
	if !exists {
		panic("User role not found in context")
	}

	role, ok := value.(enums.Role)
	if !ok {
		panic("User role is not of the correct type")
	}

	return role
}
