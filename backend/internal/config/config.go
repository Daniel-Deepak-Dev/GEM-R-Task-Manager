package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI string
	Port     string
}

// LoadConfig reads the .env file and environment variables
func LoadConfig() Config {
	// We ignore the error here because in production we might not have a .env file
	// and just use real environment variables instead.
	_ = godotenv.Load()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		MongoURI: mongoURI,
		Port:     port,
	}
}
