package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	MongoURI       string
	MongoDB        string
	JWTSecret      string
	JWTExpiryHours int
	AdminLogin     string
	AdminPassword  string
	CORSOrigins    []string
}

func Load() *Config {

	if err := godotenv.Load(); err != nil {
		log.Printf("config: .env file not loaded (%v), falling back to process env", err)
	}

	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		MongoURI:       getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:        getEnv("MONGO_DB", "exam_tasks_db"),
		JWTSecret:      strings.TrimSpace(os.Getenv("JWT_SECRET")),
		AdminLogin:     strings.TrimSpace(os.Getenv("ADMIN_LOGIN")),
		AdminPassword:  os.Getenv("ADMIN_PASSWORD"),
		CORSOrigins:    splitCSV(getEnv("CORS_ORIGINS", "*")),
		JWTExpiryHours: getEnvInt("JWT_EXPIRY_HOURS", 72),
	}

	if cfg.JWTSecret == "" {
		log.Fatal("config: JWT_SECRET is required and cannot be empty")
	}
	if cfg.AdminPassword == "" {
		log.Fatal("config: ADMIN_PASSWORD is required and cannot be empty")
	}
	if cfg.AdminLogin == "" {
		log.Fatal("config: ADMIN_LOGIN is required and cannot be empty")
	}
	if cfg.JWTExpiryHours <= 0 {
		log.Fatal(fmt.Sprintf("config: JWT_EXPIRY_HOURS must be > 0, got %d", cfg.JWTExpiryHours))
	}

	return cfg
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
		log.Printf("config: %s=%q is not a valid int, using default %d", key, v, def)
	}
	return def
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
