package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func Connect(ctx context.Context, config *Config) (*sql.DB, error) {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(config.User, config.Password),
		Host:   fmt.Sprintf("%s:%d", config.Host, config.Port),
		Path:   "/" + config.DBName,
	}
	q := u.Query()
	q.Set("sslmode", config.SSLMode)
	u.RawQuery = q.Encode()
	db, err := sql.Open("postgres", u.String())

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Apply pool settings only if configured
	if config.MaxOpenConns > 0 {
		db.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		db.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(config.ConnMaxLifetime)
	}

	if err := db.PingContext(ctx); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Failed to close database after ping failure: %v", closeErr)
		}
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
