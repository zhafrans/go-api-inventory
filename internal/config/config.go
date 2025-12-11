package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppPort string
	AppEnv  string
	
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	
	JWTSecret    string
	JWTExpireHours int
}

func LoadConfig() *Config {
	return &Config{
		AppPort: getEnv("APP_PORT", ":3000"),
		AppEnv:  getEnv("APP_ENV", "development"),
		
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "inventory_user"),
		DBPassword: getEnv("DB_PASSWORD", "password123"),
		DBName:     getEnv("DB_NAME", "inventory_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		
		JWTSecret:     getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
		JWTExpireHours: getEnvAsInt("JWT_EXPIRE_HOURS", 24),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}