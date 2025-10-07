package handlers

import (
	"github.com/muhammadolammi/rentradar/internal/database"
)

// User model helpers
func DbUserToModelsUser(dbUser database.User) User {
	return User{
		ID:          dbUser.ID,
		FirstName:   dbUser.FirstName,
		LastName:    dbUser.LastName,
		Email:       dbUser.Email,
		PhoneNumber: dbUser.PhoneNumber,
		Role:        dbUser.Role,
		CreatedAt:   dbUser.CreatedAt,
	}

}

func DbUsersToModelsUsers(dbUsers []database.User) []User {
	users := []User{}
	for _, dbUser := range dbUsers {
		users = append(users, DbUserToModelsUser(dbUser))
	}
	return users
}

// Listing model helpers

func DbListingToModelsListing(dbListing database.Listing) Listing {
	return Listing{
		ID:          dbListing.ID,
		AgentID:     dbListing.AgentID,
		Title:       dbListing.Title,
		Description: dbListing.Description,
		RentType:    dbListing.RentType,
		Price:       dbListing.Price,
		Location:    dbListing.Location,
		Latitude:    dbListing.Latitude,
		Longtitude:  dbListing.Longtitude,
		HouseType:   dbListing.HouseType,
		Verified:    dbListing.Verified,
		Images:      dbListing.Images,
		Status:      dbListing.Status,
		CreatedAt:   dbListing.CreatedAt,
	}
}

func DbListingsToModelsListings(dbListings []database.Listing) []Listing {
	listings := []Listing{}
	for _, dbListing := range dbListings {
		listings = append(listings, DbListingToModelsListing(dbListing))
	}
	return listings
}
