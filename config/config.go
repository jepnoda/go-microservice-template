package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type redisConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DB       int
}

type httpConfig struct {
	Port int
	Host string
}

func LoadConfig() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func GetRedisConfig() redisConfig {
	return redisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnvAsInt("REDIS_PORT", 6379),
		Username: getEnv("REDIS_USERNAME", ""),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       getEnvAsInt("REDIS_DB", 0),
	}
}

func GetHttpConfig() httpConfig {
	return httpConfig{
		Host: getEnv("HTTP_HOST", "0.0.0.0"),
		Port: getEnvAsInt("HTTP_PORT", 8080),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
