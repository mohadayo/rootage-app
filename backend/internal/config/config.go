package config

import (
	"log"
	"os"
)

type Config struct {
	DBURL              string
	JWTSecret          string
	OpenAIAPIKey       string
	ServerPort         string
	CorsOrigin         string
	AllowedEmailDomain string
	ResendAPIKey       string
	BaseURL            string
}

func Load() *Config {
	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		log.Println("WARNING: JWT_SECRET is not set, using dev default (unsafe for production)")
		jwtSecret = "dev-secret-key"
	}

	return &Config{
		DBURL:              getEnv("DB_URL", "postgres://ses_quiz_user:ses_quiz_pass@localhost:5432/ses_quiz?sslmode=disable"),
		JWTSecret:          jwtSecret,
		OpenAIAPIKey:       getEnv("OPENAI_API_KEY", ""),
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		CorsOrigin:         getEnv("CORS_ORIGIN", "http://localhost:5173"),
		AllowedEmailDomain: getEnv("ALLOWED_EMAIL_DOMAIN", ""),
		ResendAPIKey:       getEnv("RESEND_API_KEY", ""),
		BaseURL:            getEnv("BASE_URL", getEnv("CORS_ORIGIN", "http://localhost:5173")),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
