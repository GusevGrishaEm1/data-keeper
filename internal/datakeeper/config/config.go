package config

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config file
type Config struct {
	// Port service port
	Port int
	// HTTPS is secure http
	HTTPS bool
	// Postgres postgres config
	Postgres Postgres
	// AuthService auth service config
	AuthService AuthService
}

// Postgres postgres config
type Postgres struct {
	Host     string
	Port     int
	DB       string
	User     string
	Password string
}

// AuthService auth service config
type AuthService struct {
	Host    string
	Port    int
	Timeout time.Duration
	JWTKey  string
}

// LoadConfig load config
func LoadConfig() (*Config, error) {
	// Load .env file if exists
	_ = godotenv.Load()

	// Define flags
	port := flag.Int("port", getEnvAsInt("PORT", 8080), "service port")
	https := flag.Bool("https", getEnvAsBool("HTTPS", false), "enable HTTPS")
	postgresHost := flag.String("postgres_host", getEnv("POSTGRES_HOST", "localhost"), "Postgres host")
	postgresPort := flag.Int("postgres_port", getEnvAsInt("POSTGRES_PORT", 5432), "Postgres port")
	postgresDB := flag.String("postgres_db", getEnv("POSTGRES_DB", "postgres"), "Postgres database")
	postgresUser := flag.String("postgres_user", getEnv("POSTGRES_USER", "postgres"), "Postgres user")
	postgresPassword := flag.String("postgres_password", getEnv("POSTGRES_PASSWORD", ""), "Postgres password")
	authHost := flag.String("auth_host", getEnv("AUTH_HOST", "localhost"), "Auth host")
	authPort := flag.Int("auth_port", getEnvAsInt("AUTH_PORT", 50051), "Auth port")
	authTimeout := flag.Duration("auth_timeout", getEnvAsDuration("AUTH_TIMEOUT", 30*time.Second), "Auth service timeout")
	authJWTKey := flag.String("auth_jwt_key", getEnv("AUTH_JWT_KEY", ""), "Auth service JWT key")

	// Parse flags
	flag.Parse()

	// Fill config
	config := &Config{
		Port:  *port,
		HTTPS: *https,
		Postgres: Postgres{
			Host:     *postgresHost,
			Port:     *postgresPort,
			DB:       *postgresDB,
			User:     *postgresUser,
			Password: *postgresPassword,
		},
		AuthService: AuthService{
			Host:    *authHost,
			Port:    *authPort,
			Timeout: *authTimeout,
			JWTKey:  *authJWTKey,
		},
	}

	return config, nil
}

// getEnv gets the environment variable value or returns a default value
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets the environment variable value as int or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.ParseInt(valStr, 0, 16); err == nil {
		return int(val)
	}
	return defaultValue
}

// getEnvAsBool gets the environment variable value as bool or returns a default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}

// getEnvAsDuration gets the environment variable value as duration or returns a default value
func getEnvAsDuration(name string, defaultVal time.Duration) time.Duration {
	valStr := getEnv(name, "")
	if val, err := time.ParseDuration(valStr); err == nil {
		return val
	}
	return defaultVal
}
