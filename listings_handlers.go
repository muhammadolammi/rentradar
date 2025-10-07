package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/muhammadolammi/rentradar/internal/database"
)

func (apiConfig *Config) getListingsHandler(w http.ResponseWriter, r *http.Request) {
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
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error decoding request body. err: %v", err))
		return
	}
	//  Manage nullable parameters
	location_param := sql.NullString{Valid: false}
	min_price_param := sql.NullInt64{Valid: false}
	max_price_param := sql.NullInt64{Valid: false}
	type_param := sql.NullString{Valid: false}
	if location != "" {
		location_param = sql.NullString{Valid: true, String: location}
	}
	if list_type != "" {
		type_param = sql.NullString{Valid: true, String: list_type}
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
		Offset:   int32(offset),
		Limit:    int32(body.Limit),
		Location: location_param,
		MinPrice: min_price_param,
		MaxPrice: max_price_param,
		Type:     type_param,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error geting listings. err: %v", err))
		return
	}

	converted_listings := dbListingsToListings(listings)
	respondWithJson(w, http.StatusOK, converted_listings)

}

func (apiConfig *Config) postListingsHandler(w http.ResponseWriter, r *http.Request) {
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
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error decoding request body. err: %v", err))
		return
	}
	if body.Title == "" {
		respondWithError(w, http.StatusInternalServerError, "Enter the listing title.")
		return
	}
	if body.Description == "" {
		respondWithError(w, http.StatusInternalServerError, "Enter the listing description.")
		return
	}
	if body.HouseType == "" {
		respondWithError(w, http.StatusInternalServerError, "Enter the listing house_type.")
		return
	}
	if body.RentType == "" {
		respondWithError(w, http.StatusInternalServerError, "Enter the listing rent_type.")
		return
	}
	if len(body.Images) == 0 || string(body.Images) == "[]" || string(body.Images) == "{}" {
		respondWithError(w, http.StatusInternalServerError, "Enter the listing images.")
		return
	}
	if body.Price == 0 {
		respondWithError(w, http.StatusInternalServerError, "Enter the listing price.")
		return
	}
	if body.Location == "" {
		respondWithError(w, http.StatusInternalServerError, "Enter the listing location.")
		return
	}
	// TODO get agent id, we can handle this endpoints using midleware that makes sure this endpoint can only be called by an agent then inject the agent payload to the request body.
	listing, err := apiConfig.DB.CreateListing(context.Background(), database.CreateListingParams{
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
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating listing. err: %v", err))
		return
	}
	respondWithJson(w, http.StatusOK, dbListingToListing(listing))
}
