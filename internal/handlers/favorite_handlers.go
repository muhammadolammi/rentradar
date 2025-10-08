package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/muhammadolammi/rentradar/internal/database"
	"github.com/muhammadolammi/rentradar/internal/helpers"
)

// ---------- Create Favorite ----------
func (apiConfig *Config) PostFavoritesHandler(w http.ResponseWriter, r *http.Request, user User) {

	body := struct {
		ListingID uuid.UUID `json:"listing_id"`
	}{}

	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	// Validate input
	if body.ListingID == uuid.Nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Enter the listing_id.")
		return
	}

	// Create favorite
	fav, err := apiConfig.DB.CreateFavorite(context.Background(), database.CreateFavoriteParams{
		UserID:    user.ID,
		ListingID: body.ListingID,
	})
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "error saving favorite. err: "+err.Error())
		return
	}

	helpers.RespondWithJson(w, http.StatusOK, DbFavoriteToModelFavorite(fav))
}

// ---------- Get User Favorites ----------
func (apiConfig *Config) GetFavoritesHandler(w http.ResponseWriter, r *http.Request, user User) {

	favs, err := apiConfig.DB.GetUserFavorites(context.Background(), user.ID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "error fetching favorites. err: "+err.Error())
		return
	}

	helpers.RespondWithJson(w, http.StatusOK, DbFavoritesToModelFavorites(favs))
}
