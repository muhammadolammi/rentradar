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
		ID:           dbListing.ID,
		AgentID:      dbListing.AgentID,
		Title:        dbListing.Title,
		Description:  dbListing.Description,
		Price:        dbListing.Price,
		Location:     dbListing.Location,
		Latitude:     dbListing.Latitude,
		Longtitude:   dbListing.Longtitude,
		PropertyType: dbListing.PropertyType,
		Verified:     dbListing.Verified,
		Images:       dbListing.Images,
		Status:       dbListing.Status,
		CreatedAt:    dbListing.CreatedAt,
	}
}

func DbListingsToModelsListings(dbListings []database.Listing) []Listing {
	listings := []Listing{}
	for _, dbListing := range dbListings {
		listings = append(listings, DbListingToModelsListing(dbListing))
	}
	return listings
}

// Alert Model Helper
func DbAlertToModelsAlert(dbAlert database.Alert) Alert {
	return Alert{
		ID:            dbAlert.ID,
		UserID:        dbAlert.UserID,
		MinPrice:      dbAlert.MinPrice,
		MaxPrice:      dbAlert.MaxPrice,
		Location:      dbAlert.Location,
		PropertyType:  dbAlert.PropertyType,
		ContactMethod: dbAlert.ContactMethod,
	}
}

func DbAlertsToModelsAlerts(dbAlerts []database.Alert) []Alert {
	alerts := []Alert{}
	for _, dbAlert := range dbAlerts {
		alerts = append(alerts, DbAlertToModelsAlert(dbAlert))
	}
	return alerts
}

// Favourite Model helper
func DbFavoriteToModelFavorite(dbFav database.Favorite) Favorite {
	return Favorite{
		ID:        dbFav.ID,
		UserID:    dbFav.UserID,
		ListingID: dbFav.ListingID,
	}
}

func DbFavoritesToModelFavorites(dbFavs []database.Favorite) []Favorite {
	favs := []Favorite{}
	for _, f := range dbFavs {
		favs = append(favs, DbFavoriteToModelFavorite(f))
	}
	return favs
}

// Notification  Model Helper
func DbNotificationToModelsNotification(dbNotification database.Notification) Notification {
	return Notification{
		ID:            dbNotification.ID,
		UserID:        dbNotification.UserID,
		Status:        dbNotification.Status,
		ListingID:     dbNotification.ListingID,
		SentAt:        dbNotification.SentAt,
		ContactMethod: dbNotification.ContactMethod,
		Contact:       dbNotification.Contact,
		Subject:       dbNotification.Subject,
		Body:          dbNotification.Body,
	}
}

func DbNotificationsToModelsNotifications(dbNotifications []database.Notification) []Notification {
	notications := []Notification{}
	for _, dbNotification := range dbNotifications {
		notications = append(notications, DbNotificationToModelsNotification(dbNotification))
	}
	return notications
}
