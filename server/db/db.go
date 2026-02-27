package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	if dbURL == "" {
		dbURL = "postgres://nomad:nomad@localhost:5432/nomad_c2?sslmode=disable"
	}

	var err error
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Printf("Warning: Could not connect to database: %v. Database functionality might be limited.", err)
	} else {
		log.Println("Successfully connected to PostgreSQL")
	}

	createTables()
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS agents (
			id TEXT PRIMARY KEY,
			hostname TEXT,
			ip TEXT,
			os TEXT,
			country_code TEXT,
			last_seen TIMESTAMP,
			status TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS commands (
			id SERIAL PRIMARY KEY,
			agent_id TEXT REFERENCES agents(id),
			command TEXT,
			response TEXT,
			status TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, q := range queries {
		_, err := DB.Exec(q)
		if err != nil {
			log.Fatalf("Error creating table: %v", err)
		}
	}
}
