package persistence_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Application-drop-up/Travellle/internal/db"
	_ "github.com/lib/pq"
)

// newTestDB connects to the local Postgres instance used for development
// (see infra/docker-compose.yml) and ensures migrations are applied.
// It requires the db service from infra/docker-compose.yml to be running;
// tests are skipped if a connection cannot be established.
func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	host := envOr("DB_HOST", "localhost")
	port := envOr("DB_PORT", "5432")
	user := envOr("DB_USER", "postgres")
	password := envOr("DB_PASSWORD", "password")
	name := envOr("DB_NAME", "travellle")

	dsn := "host=" + host + " port=" + port + " user=" + user +
		" password=" + password + " dbname=" + name + " sslmode=disable"

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	if err := conn.Ping(); err != nil {
		t.Skipf("skipping: local Postgres not reachable (%v) — start it with `docker compose -f infra/docker-compose.yml up -d db`", err)
	}

	if err := db.RunMigrations(conn, migrationsDir(t)); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return conn
}

func migrationsDir(t *testing.T) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to determine caller for migrations directory lookup")
	}
	return filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "migrations")
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
