package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Application-drop-up/Travellle/internal/db"
	"github.com/Application-drop-up/Travellle/internal/router"
)

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

	r := router.New(conn, apiKey)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
