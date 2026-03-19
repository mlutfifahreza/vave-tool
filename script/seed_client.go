package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "vave_db")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	username := "admin"
	password := "admin123"
	name := "Admin Client"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	query := `
		INSERT INTO clients (name, username, password, is_active)
		VALUES ($1, $2, $3, true)
		ON CONFLICT (username) DO UPDATE SET
			name = EXCLUDED.name,
			password = EXCLUDED.password,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id
	`

	var clientID string
	err = db.QueryRow(query, name, username, string(hashedPassword)).Scan(&clientID)
	if err != nil {
		log.Fatalf("Failed to insert client: %v", err)
	}

	fmt.Printf("Client created/updated successfully!\n")
	fmt.Printf("ID: %s\n", clientID)
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Password: %s\n", password)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
