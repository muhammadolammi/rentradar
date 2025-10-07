package tests

import (
	"database/sql"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/muhammadolammi/rentradar/internal/database"
	"github.com/muhammadolammi/rentradar/internal/handlers"
)

type TestEnv struct {
	App *handlers.Config
	DB  *database.Queries
}

func SetupTestEnv(t *testing.T) *TestEnv {
	t.Helper()

	// Load environment variables from .env
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("failed to load .env: %v", err)
	}

	dbURL := os.Getenv("DB_URL")
	jwtKey := os.Getenv("JWT_KEY")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("cannot connect to test DB: %v", err)
	}

	queries := database.New(db)

	app := &handlers.Config{
		DB:     queries,
		JWTKEY: jwtKey,
	}

	return &TestEnv{
		App: app,
		DB:  queries,
	}
}
