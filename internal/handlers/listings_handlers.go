package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/muhammadolammi/rentradar/internal/database"
	"github.com/muhammadolammi/rentradar/internal/helpers"
)

func (apiConfig *Config) GetListingsHandler(w http.ResponseWriter, r *http.Request) {
	//  filters should be gootten with url param (location, type)
	location := r.URL.Query().Get("location")
	propertyType := r.URL.Query().Get("property_type_name")
	minPrice := r.URL.Query().Get("min_price")
	maxPrice := r.URL.Query().Get("max_price")
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	//  Manage nullable parameters
	locationParam := sql.NullString{Valid: false}
	minPriceParam := sql.NullInt64{Valid: false}
	maxPriceParam := sql.NullInt64{Valid: false}
	propertyTypeParam := sql.NullString{
		Valid: false,
	}
	if location != "" {
		locationParam = sql.NullString{Valid: true, String: location}
	}
	if propertyType != "" {

		propertyTypeParam = sql.NullString{
			Valid:  true,
			String: propertyType,
		}
	}
	if minPrice != "" {
		minPriceInt, err := strconv.Atoi(minPrice)
		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error converting minprice string to int. err: %v", err))
			return
		}

		minPriceParam = sql.NullInt64{Valid: true, Int64: int64(minPriceInt)}
	}
	if maxPrice != "" {
		maxPriceInt, err := strconv.Atoi(maxPrice)
		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error converting maxprice string to int. err: %v", err))
			return
		}

		maxPriceParam = sql.NullInt64{Valid: true, Int64: int64(maxPriceInt)}
	}
	offset := 0
	if page != "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error converting page string to int. err: %v", err))
			return
		}
		//  calculate Offset using the limits listings per page

		offset = (pageInt - 1) * 20
	}
	limitInt := 20
	if limit != "" {
		limitIntt, err := strconv.Atoi(limit)
		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error converting page string to int. err: %v", err))
			return
		}
		//  calculate Offset using the limits listings per page

		limitInt = limitIntt
	}

	listings, err := apiConfig.DB.GetListings(context.Background(), database.GetListingsParams{
		// Location: ,
		Offset:       int32(offset),
		Limit:        int32(limitInt),
		Location:     locationParam,
		MinPrice:     minPriceParam,
		MaxPrice:     maxPriceParam,
		PropertyType: propertyTypeParam,
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

	body := struct {
		Description  string          `json:"description"`
		Title        string          `json:"title"`
		PropertyType string          `json:"property_type"`
		Images       json.RawMessage `json:"images"`
		Price        int64           `json:"price"`
		Location     string          `json:"location"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
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
	if body.PropertyType == "" {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Enter the listing property_type.")
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
		AgentID:      user.ID,
		Price:        body.Price,
		Location:     body.Location,
		Description:  body.Description,
		Title:        body.Title,
		PropertyType: body.PropertyType,
		Images:       body.Images,
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

	id := chi.URLParam(r, "ID")
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
