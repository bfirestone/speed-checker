package database

import (
	"context"
	"fmt"

	"github.com/bfirestone/speed-checker/ent"
	"github.com/bfirestone/speed-checker/internal/config"
	_ "github.com/lib/pq"
)

// InitializeDatabase creates and initializes the database with proper schema
func InitializeDatabase(dbConfig config.DatabaseConfig) (*ent.Client, error) {
	// Open database connection
	client, err := ent.Open("postgres", dbConfig.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Create/update schema
	if err := client.Schema.Create(context.Background()); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to create database schema: %w", err)
	}

	fmt.Println("Database initialized successfully")

	return client, nil
}
