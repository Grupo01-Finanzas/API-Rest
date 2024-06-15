package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config struct to hold all configuration values
type Config struct {
	DB         *gorm.DB
	JwtSecret  string
	ServerPort string
	ServerHost string
}

// LoadConfig loads configuration from environment variables or .env file
func LoadConfig() (*Config, error) {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file, using system environment variables:", err)
	}

	// Database configuration
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	// JWT secret key
	jwtSecret := os.Getenv("JWT_SECRET")

	// Server port
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080" // Default port
	}

	// Server host
	serverHost := os.Getenv("SERVER_HOST")
	if serverHost == "" {
		serverHost = "localhost" // Default host
	}

	// Database connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPass, dbName, dbSSLMode)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Create Config struct
	cfg := &Config{
		DB:         db,
		JwtSecret:  jwtSecret,
		ServerPort: serverPort,
		ServerHost: serverHost,
	}

	return cfg, nil
}
