package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	MongoDB  MongoDBConfig
	JWT      JWTConfig
	Log      LogConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type MongoDBConfig struct {
	URI      string
	Database string
}

type JWTConfig struct {
	Secret      string
	ExpiredTime time.Duration
}

type LogConfig struct {
	Level string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env.local", ".env"); err != nil {
		// Not critical if .env doesn't exist
	}

	expiredTime, _ := time.ParseDuration(getEnv("JWT_EXPIRED_TIME", "24h"))

	return &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "MyApp"),
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
		},
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			Database: getEnv("POSTGRES_DB", "blink_db"),
			SSLMode:  getEnv("POSTGRES_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       0,
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGODB_DATABASE", "myapp_logs"),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "your-secret-key"),
			ExpiredTime: expiredTime,
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
