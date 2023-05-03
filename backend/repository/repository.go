package repository

import (
	"fmt"
	"os"
	"time"
)

type Timestamps struct {
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

type PostgresConfig struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
	SslMode  string
}

// PostgresConfigFromEnv uses the following environment variables:
// - POSTGRES_HOST
// - POSTGRES_PORT
// - POSTGRES_DB
// - POSTGRES_USER
// - POSTGRES_PASSWORD
// - POSTGRES_SSLMODE
func PostgresConfigFromEnv() PostgresConfig {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	db := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	sslMode := os.Getenv("POSTGRES_SSLMODE")

	if host == "" {
		host = "localhost"
	}

	if port == "" {
		port = "5432"
	}

	if db == "" {
		db = "dev"
	}

	if user == "" {
		user = "dev"
	}

	if password == "" {
		password = "dev"
	}

	if sslMode == "" {
		sslMode = "disable"
	}

	return PostgresConfig{
		Host:     host,
		Port:     port,
		DbName:   db,
		User:     user,
		Password: password,
		SslMode:  sslMode,
	}
}

func (pc PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", pc.Host, pc.Port, pc.User, pc.Password, pc.DbName, pc.SslMode)
}
