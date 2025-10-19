package database

import (
	"context"
	"database/sql"
	"fmt"
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
	u := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		url.QueryEscape(config.User),
		url.QueryEscape(config.Password),
		config.Host, config.Port, config.DBName, config.SSLMode)
	db, err := sql.Open("postgres", u)
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
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
