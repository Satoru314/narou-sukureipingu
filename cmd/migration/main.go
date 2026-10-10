package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Get database connection info from environment
	host := getEnv("POSTGRES_HOST", "db")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "narou")
	password := getEnv("POSTGRES_PASSWORD", "narou_password")
	dbname := getEnv("POSTGRES_DB", "narou_db")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database successfully")

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully")
}

func runMigrations(db *sql.DB) error {
	// Create novels table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS novels (
			ncode VARCHAR(20) PRIMARY KEY,
			title TEXT NOT NULL,
			author TEXT NOT NULL,
			author_url TEXT,
			description TEXT,
			novel_url TEXT NOT NULL,
			episode_count INT NOT NULL DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create novels table: %w", err)
	}
	log.Println("✓ novels table created")

	// Create episodes table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS episodes (
			ncode VARCHAR(20) NOT NULL,
			episode_number INT NOT NULL,
			subtitle TEXT NOT NULL,
			url TEXT NOT NULL,
			content TEXT,
			posted_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (ncode, episode_number),
			FOREIGN KEY (ncode) REFERENCES novels(ncode) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create episodes table: %w", err)
	}
	log.Println("✓ episodes table created")

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
