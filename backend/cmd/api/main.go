package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Application-drop-up/Travellle/internal/db"
	"github.com/Application-drop-up/Travellle/internal/router"
)

const defaultAllowedOrigins = "http://localhost:3000,http://localhost:3001"

func main() {
	conn, err := db.NewConnection()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer conn.Close()

	if err := db.RunMigrations(conn, "migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	apiKey := os.Getenv("GOOGLE_PLACES_API_KEY")
	if apiKey == "" {
		log.Fatal("GOOGLE_PLACES_API_KEY is not set")
	}

	origins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if origins == "" {
		origins = defaultAllowedOrigins
	}

	r := router.New(conn, apiKey, strings.Split(origins, ","))

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
