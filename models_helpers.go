package main

import "github.com/muhammadolammi/rentradar/internal/database"

// User model helpers
func dbUserToUser(dbUser database.User) User {
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

func dbUsersToUsers(dbUsers []database.User) []User {
	users := []User{}
	for _, dbUser := range dbUsers {
		users = append(users, dbUserToUser(dbUser))
	}
	return users
}

// Listing model helpers

func dbListingToListing(dbListing database.Listing) Listing {
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

func dbListingsToListings(dbListings []database.Listing) []Listing {
	listings := []Listing{}
	for _, dbListing := range dbListings {
		listings = append(listings, dbListingToListing(dbListing))
	}
	return listings
}
