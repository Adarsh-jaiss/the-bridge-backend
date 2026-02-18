package config

import (
	"os"

	"github.com/adarsh-jaiss/the-bridge/db"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string
	Port string
	JWTSecret string
	DBConfig db.Config
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		value = fallback
	}

	return value
}

func Load() *Config {
	// for loading enviroment (ENV=development go run .)
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	err := godotenv.Load(".env" + env); if err!= nil {
		// TODO : Add logging
		panic(err)
	}
	return &Config{
		AppEnv: getEnv("APP_ENV",env),
		Port: getEnv("PORT","8080"),
		JWTSecret: getEnv("JWT_SECRET",""),
		DBConfig: db.Config{
			Host: getEnv("HOST",""),
			Port: getEnv("PORT",""),
			DBName: getEnv("DBNAME",""),
			Password: getEnv("PASSWORD",""),
			User: getEnv("USER",""),
		},
	}

}