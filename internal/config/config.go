package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port string

	// Database
	MongoURI string
	DBName   string

	// JWT
	JWTSecret     string
	JWTExpiration string

	// Environment
	Env string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	envFileName := fmt.Sprintf(".env.%s", env)
	// Try multiple paths so it works from project root or from cmd/
	envPaths := []string{
		filepath.Join("config", "environments", envFileName),
		filepath.Join("..", "config", "environments", envFileName),
	}
	var loaded bool
	for _, envFile := range envPaths {
		if err := godotenv.Load(envFile); err == nil {
			loaded = true
			break
		}
	}
	if !loaded {
		log.Printf("Warning: Could not load %s from config/environments/. Try running from project root or create the file.", envFileName)
	}

	config := &Config{
		Port:          getEnv("PORT"),
		MongoURI:      getEnv("MONGO_URI"),
		DBName:        getEnv("MONGO_DB"),
		JWTSecret:     getEnv("JWT_SECRET"),
		JWTExpiration: getEnv("JWT_EXPIRATION"),
		Env:           getEnv("APP_ENV"),
	}

	// Validate required fields
	if config.MongoURI == "" {
		log.Fatal("MONGO_URI environment variable is required")
	}

	if config.DBName == "" {
		log.Fatal("MONGO_DB environment variable is required")
	}

	if config.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	return config
}

// getEnv gets an environment variable with a default value
func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return ""
}
