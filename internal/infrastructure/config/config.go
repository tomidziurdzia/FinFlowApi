package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
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
	DatabaseURL string // Railway/Heroku style: postgres://user:pass@host:port/dbname
	Host        string
	Port        string
	User        string
	Password    string
	DBName      string
	SSLMode     string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	
	var dbConfig DatabaseConfig
	if databaseURL != "" {
		dbConfig = DatabaseConfig{
			DatabaseURL: databaseURL,
		}
	} else {
		dbConfig = DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "finflow"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		}
		
		if dbConfig.Password == "" {
			return nil, fmt.Errorf("DB_PASSWORD is required when DATABASE_URL is not set")
		}
		if dbConfig.DBName == "" {
			return nil, fmt.Errorf("DB_NAME is required when DATABASE_URL is not set")
		}
	}

	cfg := &Config{
		Port:     getEnv("PORT", "8080"),
		Database: dbConfig,
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

	return cfg, nil
}

func (c *DatabaseConfig) ConnectionString() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func (c *DatabaseConfig) ParseDatabaseURL() error {
	if c.DatabaseURL == "" {
		return nil
	}

	parsedURL, err := url.Parse(c.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse DATABASE_URL: %w", err)
	}

	c.User = parsedURL.User.Username()
	c.Password, _ = parsedURL.User.Password()
	
	hostPort := strings.Split(parsedURL.Host, ":")
	c.Host = hostPort[0]
	if len(hostPort) > 1 {
		c.Port = hostPort[1]
	} else {
		c.Port = "5432"
	}

	c.DBName = strings.TrimPrefix(parsedURL.Path, "/")
	
	if parsedURL.Query().Get("sslmode") != "" {
		c.SSLMode = parsedURL.Query().Get("sslmode")
	} else {
		c.SSLMode = "require"
	}

	return nil
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
