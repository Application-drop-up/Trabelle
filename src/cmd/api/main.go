package main

import (
	"log"
	"net/http"

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

	r := router.New(conn)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
