package config

import "os"

type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	FrontendOrigin string
}

func Load() Config {
	return Config{
		Port:           getenv("PORT", "8080"),
		DatabaseURL:    getenv("DATABASE_URL", "postgres://mohitrawat@localhost:5432/istream?sslmode=disable"),
		JWTSecret:      getenv("JWT_SECRET", "dev-secret-change-in-production"),
		FrontendOrigin: getenv("FRONTEND_ORIGIN", "http://localhost:5173"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
