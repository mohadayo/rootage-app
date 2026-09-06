package config

import (
	"log"
	"os"
)

type Config struct {
	AppEnv             string
	DBURL              string
	JWTSecret          string
	OpenAIAPIKey       string
	ServerPort         string
	CorsOrigin         string
	AllowedEmailDomain string
	ResendAPIKey       string
	MailFrom           string
	BaseURL            string
	AdminEmail         string
	AdminPassword      string
}

func Load() *Config {
	appEnv := getEnv("APP_ENV", "production")

	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		// 本番では「動くが安全でない既定値」を許さず起動を止める。
		// 開発時のみ既定値でのフォールバックを許容する。
		if appEnv != "development" {
			log.Fatalf("JWT_SECRET is not set (APP_ENV=%s). Set a strong secret or run with APP_ENV=development.", appEnv)
		}
		log.Println("WARNING: JWT_SECRET is not set, using dev default (APP_ENV=development only)")
		jwtSecret = "dev-secret-key"
	}

	return &Config{
		AppEnv:             appEnv,
		DBURL:              getEnv("DB_URL", "postgres://ses_quiz_user:ses_quiz_pass@localhost:5432/ses_quiz?sslmode=disable"),
		JWTSecret:          jwtSecret,
		OpenAIAPIKey:       getEnv("OPENAI_API_KEY", ""),
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		CorsOrigin:         getEnv("CORS_ORIGIN", "http://localhost:5173"),
		AllowedEmailDomain: getEnv("ALLOWED_EMAIL_DOMAIN", ""),
		ResendAPIKey:       getEnv("RESEND_API_KEY", ""),
		// リセットメールの送信元。onboarding@resend.dev は Resend のサンドボックス用で
		// アカウント所有者本人にしか配送されないため、本番では検証済みドメインの
		// アドレスを MAIL_FROM で必ず設定すること。
		MailFrom:      getEnv("MAIL_FROM", "rootage <onboarding@resend.dev>"),
		BaseURL:       getEnv("BASE_URL", getEnv("CORS_ORIGIN", "http://localhost:5173")),
		AdminEmail:    getEnv("ADMIN_EMAIL", ""),
		AdminPassword: getEnv("ADMIN_PASSWORD", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
