package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/muhammadolammi/rentradar/internal/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("there is no port provided kindly provide a port.")
		return
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Println("empty dbURL")
		return
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	dbQueries := database.New(db)

	apiConfig := Config{
		PORT: port,
		DB:   dbQueries,
	}

	server(&apiConfig)

}
