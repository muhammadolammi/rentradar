package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/muhammadolammi/rentradar/internal/handlers"
)

func server(apiConfig *handlers.Config) {

	// Define CORS options
	corsOptions := cors.Options{
		AllowedOrigins: []string{"http://localhost"},

		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // You can customize this based on your needs
		AllowCredentials: true,
		MaxAge:           300, // Maximum age for cache, in seconds
	}
	router := chi.NewRouter()
	apiRoute := chi.NewRouter()
	// ADD MIDDLREWARE
	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(cors.Handler(corsOptions))
	router.Use(apiConfig.VerifyApiKey())

	// ADD ROUTES
	apiRoute.Get("/hello", handlers.SuccessResponse)
	apiRoute.Get("/error", handlers.ErrorResponse)

	// Handle Auth
	apiRoute.Post("/register", apiConfig.RegisterHandler)
	apiRoute.Post("/login", apiConfig.LoginHandler)

	// users Handlers
	apiRoute.Get("/user", apiConfig.AuthMiddleware(false, []byte(apiConfig.JWTKEY), apiConfig.GetUserHandler))

	//  Listings handlers
	apiRoute.Get("/listings", apiConfig.GetListingsHandler)
	apiRoute.Post("/listings", apiConfig.AuthMiddleware(false, []byte(apiConfig.JWTKEY), apiConfig.PostListingsHandler))
	apiRoute.Get("/listings/{ID}", apiConfig.GetListingHandler)

	// Listing handlers
	apiRoute.Post("/listings", apiConfig.AuthMiddleware(false, []byte(apiConfig.JWTKEY), apiConfig.PostListingsHandler))
	router.Get("/listings/{ID}", apiConfig.GetListingHandler)
	router.Get("/listings", apiConfig.GetListingsHandler)

	// alert handlers
	router.Post("/alerts", apiConfig.AuthMiddleware(false, []byte(apiConfig.JWTKEY), apiConfig.PostAlertsHandler))
	router.Get("/alerts", apiConfig.AuthMiddleware(false, []byte(apiConfig.JWTKEY), apiConfig.GetAlertsHandler))

	// favorite handlers
	router.Post("/favorites", apiConfig.AuthMiddleware(false, []byte(apiConfig.JWTKEY), apiConfig.PostFavoritesHandler))
	router.Get("/favorites", apiConfig.AuthMiddleware(false, []byte(apiConfig.JWTKEY), apiConfig.GetFavoritesHandler))

	router.Mount("/api", apiRoute)
	srv := &http.Server{
		Addr:              ":" + apiConfig.PORT,
		Handler:           router,
		ReadHeaderTimeout: time.Minute,
	}

	log.Printf("Serving on port: %s\n", apiConfig.PORT)
	log.Fatal(srv.ListenAndServe())
}
