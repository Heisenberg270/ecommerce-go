package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func initDB() {
	// build DSN from env
	dsn := fmt.Sprintf(
		"host=db port=5432 user=postgres password=%s dbname=ecommerce sslmode=disable",
		os.Getenv("POSTGRES_PASSWORD"),
	)

	var err error
	// retry loop: try up to 10 times, sleeping 2s between
	for i := 1; i <= 10; i++ {
		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			log.Println("Connected to DB")
			break
		}
		log.Printf("DB not ready (attempt %d): %v", i, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Unable to connect to DB after retries: %v", err)
	}

	// migration
	schema := `
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		price NUMERIC(10,2) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
  	);`

	// users table for authentication
	schema += `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
	);`

	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("Failed to migrate DB: %v", err)
	}
}
