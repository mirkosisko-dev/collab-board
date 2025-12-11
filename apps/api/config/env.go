package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl string
}

var Envs = initConfig()

func initConfig() Config {
	loadDotEnv()

	return Config{
		DatabaseUrl: getEnv("DATABASE_URL", "postgres://admin:admin@localhost:5432/collab_board?sslmode=disable"),
	}
}

// loadDotEnv walks up from the current working directory to locate a .env file.
// This covers running the binary from nested folders (e.g., apps/api/cmd/...).
func loadDotEnv() {
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	// Try up to 6 levels up to be safe in nested dirs
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
