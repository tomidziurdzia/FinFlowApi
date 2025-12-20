package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_DefaultValues(t *testing.T) {
	originalEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, env := range originalEnv {
			key := env[:len(env)-len(os.Getenv(env))-1]
			os.Setenv(key, os.Getenv(key))
		}
	}()

	os.Clearenv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Port != "8080" {
		t.Errorf("expected default port 8080, got %s", cfg.Port)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("expected default host localhost, got %s", cfg.Database.Host)
	}

	if cfg.Database.Port != "5432" {
		t.Errorf("expected default port 5432, got %s", cfg.Database.Port)
	}

	if cfg.App.SystemUser != "system" {
		t.Errorf("expected default system user 'system', got %s", cfg.App.SystemUser)
	}
}

func TestLoad_WithEnvVars(t *testing.T) {
	originalEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, env := range originalEnv {
			key := env[:len(env)-len(os.Getenv(env))-1]
			os.Setenv(key, os.Getenv(key))
		}
	}()

	os.Clearenv()
	os.Setenv("PORT", "3000")
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("APP_SYSTEM_USER", "admin")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Port != "3000" {
		t.Errorf("expected port 3000, got %s", cfg.Port)
	}

	if cfg.Database.Host != "testhost" {
		t.Errorf("expected host testhost, got %s", cfg.Database.Host)
	}

	if cfg.Database.Port != "5433" {
		t.Errorf("expected port 5433, got %s", cfg.Database.Port)
	}

	if cfg.App.SystemUser != "admin" {
		t.Errorf("expected system user 'admin', got %s", cfg.App.SystemUser)
	}
}

func TestLoad_WithDatabaseURL(t *testing.T) {
	originalEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, env := range originalEnv {
			key := env[:len(env)-len(os.Getenv(env))-1]
			os.Setenv(key, os.Getenv(key))
		}
	}()

	os.Clearenv()
	databaseURL := "postgres://user:pass@host:5432/dbname?sslmode=require"
	os.Setenv("DATABASE_URL", databaseURL)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Database.DatabaseURL != databaseURL {
		t.Errorf("expected DATABASE_URL %s, got %s", databaseURL, cfg.Database.DatabaseURL)
	}
}

func TestConnectionString_WithDatabaseURL(t *testing.T) {
	cfg := DatabaseConfig{
		DatabaseURL: "postgres://user:pass@host:5432/dbname",
	}

	connStr := cfg.ConnectionString()
	if connStr != cfg.DatabaseURL {
		t.Errorf("expected connection string to be DATABASE_URL, got %s", connStr)
	}
}

func TestConnectionString_WithoutDatabaseURL(t *testing.T) {
	cfg := DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "password",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	connStr := cfg.ConnectionString()
	expected := "host=localhost port=5432 user=postgres password=password dbname=testdb sslmode=disable"

	if connStr != expected {
		t.Errorf("expected connection string %s, got %s", expected, connStr)
	}
}

func TestServerConfig_Timeouts(t *testing.T) {
	originalEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, env := range originalEnv {
			key := env[:len(env)-len(os.Getenv(env))-1]
			os.Setenv(key, os.Getenv(key))
		}
	}()

	os.Clearenv()
	os.Setenv("SERVER_READ_TIMEOUT", "20")
	os.Setenv("SERVER_WRITE_TIMEOUT", "40")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Server.ReadTimeout != 20*time.Second {
		t.Errorf("expected ReadTimeout 20s, got %v", cfg.Server.ReadTimeout)
	}

	if cfg.Server.WriteTimeout != 40*time.Second {
		t.Errorf("expected WriteTimeout 40s, got %v", cfg.Server.WriteTimeout)
	}
}