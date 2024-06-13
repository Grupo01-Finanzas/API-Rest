package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func CorsMiddleware() gin.HandlerFunc {
	// Configure CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Accept"},
		AllowCredentials: true,
	})

	// Wrap the cors.HandlerFunc with a gin.HandlerFunc
	return func(ctx *gin.Context) {
		c.HandlerFunc(ctx.Writer, ctx.Request)
		ctx.Next()
	}
}
