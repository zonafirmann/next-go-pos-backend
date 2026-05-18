package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

// ConnectDB reads the environment variables and establishes a connection to PostgreSQL.
func ConnectDB() *pgx.Conn {
	// 1. Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Warning: .env file not found, falling back to system environment variables.")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("❌ FATAL ERROR: DATABASE_URL is not set in the environment.")
	}

	// 2. Establish connection using the pgx driver
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("❌ FATAL ERROR: Unable to connect to PostgreSQL database: %v\n", err)
	}

	fmt.Println("✅ [DATABASE] Enterprise PostgreSQL connection established successfully!")
	return conn
}
