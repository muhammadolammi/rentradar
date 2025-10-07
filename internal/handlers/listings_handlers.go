package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/muhammadolammi/rentradar/internal/database"
	"github.com/muhammadolammi/rentradar/internal/helpers"
)

func (apiConfig *Config) GetListingsHandler(w http.ResponseWriter, r *http.Request) {
	// things like limit and page should be gotten using body
	body := struct {
		MinPrice int64 `json:"min_price"`
		MaxPrice int64 `json:"max_price"`
		Page     int   `json:"page"`
		Limit    int   `json:"limit"`
	}{}
	//  filters should be gootten with url param (location, type)
	location := r.URL.Query().Get("location")
	list_type := r.URL.Query().Get("type")

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error decoding request body. err: %v", err))
		return
	}
	//  Manage nullable parameters
	location_param := sql.NullString{Valid: false}
	min_price_param := sql.NullInt64{Valid: false}
	max_price_param := sql.NullInt64{Valid: false}
	house_type_param := sql.NullString{Valid: false}
	if location != "" {
		location_param = sql.NullString{Valid: true, String: location}
	}
	if list_type != "" {
		house_type_param = sql.NullString{Valid: true, String: list_type}
	}
	if body.MinPrice != 0 {

		min_price_param = sql.NullInt64{Valid: true, Int64: body.MinPrice}
	}
	if body.MaxPrice != 0 {

		max_price_param = sql.NullInt64{Valid: true, Int64: body.MaxPrice}
	}

	//  calculate Offset using the limits listings per page
	offset := (body.Page - 1) * body.Limit

	listings, err := apiConfig.DB.GetListings(context.Background(), database.GetListingsParams{
		// Location: ,
		Offset:    int32(offset),
		Limit:     int32(body.Limit),
		Location:  location_param,
		MinPrice:  min_price_param,
		MaxPrice:  max_price_param,
		HouseType: house_type_param,
	})
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error geting listings. err: %v", err))
		return
	}

	converted_listings := DbListingsToModelsListings(listings)
	helpers.RespondWithJson(w, http.StatusOK, converted_listings)

}

func (apiConfig *Config) PostListingsHandler(w http.ResponseWriter, r *http.Request, user User) {
	if user.Role != "agent" {
		helpers.RespondWithError(w, http.StatusUnauthorized, "user not an agent")
		return
	}
	agent, err := apiConfig.DB.GetAgentWithUserId(r.Context(), user.ID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf("user getting agent. err: %v", err))
		return
	}
	body := struct {
		Description string          `json:"description"`
		Title       string          `json:"title"`
		RentType    string          `json:"rent_type"`
		HouseType   string          `json:"house_type"`
		Images      json.RawMessage `json:"images"`
		Price       int64           `json:"price"`
		Location    string          `json:"location"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&body)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error decoding request body. err: %v", err))
		return
	}
	if body.Title == "" {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Enter the listing title.")
		return
	}
	if body.Description == "" {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Enter the listing description.")
		return
	}
	if body.HouseType == "" {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Enter the listing house_type.")
		return
	}
	if body.RentType == "" {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Enter the listing rent_type.")
		return
	}
	if len(body.Images) == 0 || string(body.Images) == "[]" || string(body.Images) == "{}" {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Enter the listing images.")
		return
	}
	if body.Price == 0 {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Enter the listing price.")
		return
	}
	if body.Location == "" {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Enter the listing location.")
		return
	}
	listing, err := apiConfig.DB.CreateListing(context.Background(), database.CreateListingParams{
		AgentID:     agent.ID,
		Price:       body.Price,
		Location:    body.Location,
		Description: body.Description,
		Title:       body.Title,
		HouseType:   body.HouseType,
		RentType:    body.RentType,
		Images:      body.Images,
		// Status should be active on creation
		Status: "active",
	})
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating listing. err: %v", err))
		return
	}
	helpers.RespondWithJson(w, http.StatusOK, DbListingToModelsListing(listing))
}

func (apiConfig *Config) GetListingHandler(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "playlistID")
	uuidID, err := uuid.Parse(id)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error parsing uuid. err: %v", err))
		return
	}

	listing, err := apiConfig.DB.GetListing(context.Background(), uuidID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error geting listing. err: %v", err))
		return
	}

	converted_listings := DbListingToModelsListing(listing)
	helpers.RespondWithJson(w, http.StatusOK, converted_listings)

}
