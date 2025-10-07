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
	api_key := os.Getenv("API_KEY")
	if dbURL == "" {
		log.Println("empty apiKEY")
		return
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	dbQueries := database.New(db)

	apiConfig := Config{
		PORT:   port,
		DB:     dbQueries,
		APIKEY: api_key,
	}

	server(&apiConfig)

}
