package config

import (
	"os"
	"sync"

	db "github.com/adarsh-jaiss/the-bridge/internal/database"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv                string
	Port                  string
	JWTAccessTokenSecret  string
	JWTRefreshTokenSecret string
	DBConfig              db.Config
	ResendAPIKey          string
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		value = fallback
	}

	return value
}

var (
	instance *Config
	once     sync.Once
)

// Get returns the global config instance, initializing it if necessary
func Get() *Config {
	once.Do(
		func() {
			instance = load()
		})
	return instance
}

func load() *Config {
	// for loading enviroment (ENV=development go run .)
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	err := godotenv.Load(".env." + env)
	if err != nil {
		// TODO : Add logging
		panic(err)
	}
	return &Config{
		AppEnv:                getEnv("APP_ENV", env),
		Port:                  getEnv("PORT", "8080"),
		JWTAccessTokenSecret:  getEnv("JWT_ACCESS_TOKEN_SECRET", ""),
		JWTRefreshTokenSecret: getEnv("JWT_REFRESH_TOKEN_SECRET", ""),
		DBConfig: db.Config{
			Host:     getEnv("DB_HOST", ""),
			Port:     getEnv("DB_PORT", ""),
			DBName:   getEnv("DBNAME", ""),
			Password: getEnv("DB_PASSWORD", ""),
			User:     getEnv("DB_USER", ""),
			SSLMode:  getEnv("DB_SSL_MODE", ""),
		},
		ResendAPIKey: getEnv("RESEND_API_KEY", ""),
	}
}
