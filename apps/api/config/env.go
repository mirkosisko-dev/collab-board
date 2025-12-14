package config

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl                     string
	JWTExpirationInSeconds          int64
	JWTSecret                       string
	RefreshTokenExpirationInSeconds int64
	RefreshTokenSecret              string
}

var Envs = initConfig()

func initConfig() Config {
	loadDotEnv()

	return Config{
		DatabaseUrl:                     getEnv("DATABASE_URL", "postgres://admin:admin@localhost:5432/collab_board?sslmode=disable"),
		JWTExpirationInSeconds:          getEnvAsInt("JWT_EXP", 3600*24*7),
		JWTSecret:                       getEnv("JWT_SECRET", "not-a-secret-any-more"),
		RefreshTokenExpirationInSeconds: getEnvAsInt("RT_EXP", 3600*24*100),
		RefreshTokenSecret:              getEnv("RT_SECRET", "not-a-secret-any-more-2"),
	}
}

func loadDotEnv() {
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	for i := 0; i < 6; i++ {
		path := filepath.Join(wd, ".env")
		if _, err := os.Stat(path); err == nil {
			_ = godotenv.Load(path)
			return
		}
		wd = filepath.Dir(wd)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
