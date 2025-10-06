package main

import "github.com/muhammadolammi/rentradar/internal/database"

type Config struct {
	DB   *database.Queries
	PORT string
}
