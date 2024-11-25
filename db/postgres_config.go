package db

import (
	"context"
	"fmt"
	"log"
	"main/ent"
	"os"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
)

var client *ent.Client

func InitDB() {
	// Get database connection parameters from env
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")
	sslmode := os.Getenv("DB_SSLMODE")

	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host, port, user, dbname, password, sslmode)

	// Open the database connection
	drv, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	client = ent.NewClient(ent.Driver(drv))

	// Run migrations
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed to run schema migrations: %v", err)
	}

	log.Println("Database connection established.")
}

func InitDBDefault() {
	drv, err := sql.Open(dialect.Postgres, "host=localhost port=5432 user=temp dbname=assessment password=temp sslmode=disable")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	client = ent.NewClient(ent.Driver(drv))

	// Run migrations
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed to run schema migrations: %v", err)
	}

	log.Println("Database connection established.")
}

func GetClient() *ent.Client {
	if client == nil {
		log.Fatal("Database client is not initialized.")
	}
	return client
}

func CloseDB() {
	if client != nil {
		client.Close()
		log.Println("Database connection closed.")
	}
}
