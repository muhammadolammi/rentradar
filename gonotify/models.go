package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/muhammadolammi/rentradar/internal/database"
)

type SMTPModel struct {
	Server   string
	Password string
	UserName string
}

type Config struct {
	DB          *database.Queries
	SMTPModel   SMTPModel
	RABBITMQUrl string
}

type Notification struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	Contact       string    `json:"contact"`
	Body          string    `json:"body"`
	Subject       string    `json:"subject"`
	ListingID     uuid.UUID `json:"listing_id"`
	SentAt        time.Time `json:"sent_at"`
	Status        string    `json:"status"`
	ContactMethod string    `json:"contact_method"`
}
