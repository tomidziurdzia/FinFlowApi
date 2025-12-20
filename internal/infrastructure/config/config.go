package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	Database  DatabaseConfig
	Server    ServerConfig
	App       AppConfig
}

type ServerConfig struct {
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout        time.Duration
	ShutdownTimeout    time.Duration
}

type AppConfig struct {
	SystemUser string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Port: getEnv("PORT", "8080"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "finflow"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			ReadTimeout:       getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
			ReadHeaderTimeout: getDurationEnv("SERVER_READ_HEADER_TIMEOUT", 5*time.Second),
			WriteTimeout:      getDurationEnv("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:        getDurationEnv("SERVER_IDLE_TIMEOUT", 1*time.Minute),
			ShutdownTimeout:    getDurationEnv("SERVER_SHUTDOWN_TIMEOUT", 10*time.Second),
		},
		App: AppConfig{
			SystemUser: getEnv("APP_SYSTEM_USER", "system"),
		},
	}

	if cfg.Database.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD es requerida")
	}
	if cfg.Database.DBName == "" {
		return nil, fmt.Errorf("DB_NAME es requerida")
	}

	return cfg, nil
}

func (c *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if seconds, err := strconv.Atoi(value); err == nil {
			return time.Duration(seconds) * time.Second
		}
	}
	return defaultValue
}
