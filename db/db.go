package db

import (
	"database/sql"
	"fmt"
	"log"
	"oravue_backend/internal/config"

	_ "github.com/lib/pq"
)

type Postgresql struct {
	Db *sql.DB
}

// New initializes a new PostgreSQL connection and ensures the table is created
func New(cfg *config.Config) (*Postgresql, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			phone_number TEXT NOT NULL UNIQUE,
			otp TEXT NOT NULL,
			expiry TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	log.Println("Database connection established and table ensured.")

	return &Postgresql{Db: db}, nil
}
