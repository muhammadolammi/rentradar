package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/muhammadolammi/rentradar/internal/database"
	"github.com/muhammadolammi/rentradar/internal/helpers"
)

// ---------- Create Alert ----------
func (apiConfig *Config) PostAlertsHandler(w http.ResponseWriter, r *http.Request, user User) {

	body := struct {
		MinPrice      int64  `json:"min_price"`
		MaxPrice      int64  `json:"max_price"`
		Location      string `json:"location"`
		PropertyType  string `json:"property_type"`
		ContactMethod string `json:"contact_method"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if body.MinPrice == 0 {
		helpers.RespondWithError(w, http.StatusBadRequest, "Enter the min_price.")
		return
	}
	if body.MaxPrice == 0 {
		helpers.RespondWithError(w, http.StatusBadRequest, "Enter the max_price.")
		return
	}
	if body.Location == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Enter the location.")
		return
	}
	if body.ContactMethod == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Enter the contact method.")
		return
	}

	alert, err := apiConfig.DB.CreateAlert(context.Background(), database.CreateAlertParams{
		UserID:        user.ID,
		MinPrice:      body.MinPrice,
		MaxPrice:      body.MaxPrice,
		Location:      body.Location,
		PropertyType:  body.PropertyType,
		ContactMethod: body.ContactMethod,
	})

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "error creating alert. err: "+err.Error())
		return
	}

	helpers.RespondWithJson(w, http.StatusOK, alert)
}

// ---------- Get User Alerts ----------
func (apiConfig *Config) GetAlertsHandler(w http.ResponseWriter, r *http.Request, user User) {

	alerts, err := apiConfig.DB.GetUserAlerts(context.Background(), user.ID)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "error getting user alerts: "+err.Error())
		return
	}

	helpers.RespondWithJson(w, http.StatusOK, alerts)
}
