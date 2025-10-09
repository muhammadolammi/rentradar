package tests

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/muhammadolammi/rentradar/internal/database"
	"github.com/muhammadolammi/rentradar/internal/handlers"
)

type TestEnv struct {
	App    *handlers.Config
	DB     *database.Queries
	Router *chi.Mux
}

func SetupTestEnv(t *testing.T) *TestEnv {
	t.Helper()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		t.Fatalf("failed to load .env: %v", err)
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Println("empty dbURL")

	}
	api_key := os.Getenv("API_KEY")
	if dbURL == "" {
		log.Println("empty apiKEY")
	}
	jwt_key := os.Getenv("JWT_KEY")
	if jwt_key == "" {
		log.Println("empty jwtKEY")
	}
	sudo_key := os.Getenv("SUDO_KEY")
	if sudo_key == "" {
		log.Println("empty sudoKEY")
	}

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		t.Fatalf("cannot connect to test DB: %v", err)
	}

	queries := database.New(db)

	app := &handlers.Config{
		DB:      queries,
		JWTKEY:  jwt_key,
		APIKEY:  api_key,
		SUDOKEY: sudo_key,
	}

	// 🔹 Setup Chi router for tests
	router := chi.NewRouter()

	// 🔹 Use same middlewares as in your production server
	corsOptions := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}
	router.Use(cors.Handler(corsOptions))

	router.Use(app.VerifyApiKey())

	router.Post("/register", app.RegisterHandler)
	router.Post("/login", app.LoginHandler)
	router.Post("/refresh", app.RefreshTokens)

	router.Post("/listings", app.AuthMiddleware(false, []byte(jwt_key), app.PostListingsHandler))
	router.Get("/listings/{ID}", app.GetListingHandler)
	router.Get("/listings", app.GetListingsHandler)
	router.Post("/alerts", app.AuthMiddleware(false, []byte(jwt_key), app.PostAlertsHandler))
	router.Get("/alerts", app.AuthMiddleware(false, []byte(jwt_key), app.GetAlertsHandler))

	router.Post("/favorites", app.AuthMiddleware(false, []byte(jwt_key), app.PostFavoritesHandler))
	router.Get("/favorites", app.AuthMiddleware(false, []byte(jwt_key), app.GetFavoritesHandler))

	return &TestEnv{
		App:    app,
		DB:     queries,
		Router: router,
	}
}
