package handlers

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/muhammadolammi/rentradar/internal/database"
)

type Config struct {
	DB      *database.Queries
	PORT    string
	APIKEY  string
	JWTKEY  string
	SUDOKEY string
}

type Agent struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	CompanyName string    `json:"company_name"`
	Verified    bool      `json:"verified"`
	Rating      float64   `json:"rating"`
}

type Alert struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	MinPrice       int64     `json:"min_price"`
	MaxPrice       int64     `json:"max_price"`
	Location       string    `json:"location"`
	PropertyTypeID uuid.UUID `json:"property_type"`
	ContactMethod  string    `json:"contact_method"`
}

type Favorite struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ListingID uuid.UUID `json:"lsting_id"`
}

type Listing struct {
	ID             uuid.UUID       `json:"id"`
	AgentID        uuid.UUID       `json:"agent_id"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	Price          int64           `json:"price"`
	Location       string          `json:"location"`
	Latitude       sql.NullFloat64 `json:"latitude"`
	Longtitude     sql.NullFloat64 `json:"longtitude"`
	PropertyTypeId uuid.UUID       `json:"property_type_id"`
	Verified       bool            `json:"verified"`
	Images         json.RawMessage `json:"images"`
	Status         string          `json:"status"`
	CreatedAt      time.Time       `json:"created_at"`
}

type Notification struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ListingID uuid.UUID `json:"listing_id"`
	SentAt    time.Time `json:"sent_at"`
	Status    string    `json:"status"`
}

type User struct {
	ID          uuid.UUID      `json:"id"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Email       string         `json:"email"`
	PhoneNumber sql.NullString `json:"phone_number"`
	Role        string         `json:"role"`
	// Password    string         `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type PropertyType struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
