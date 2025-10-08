package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/muhammadolammi/rentradar/internal/helpers"
)

func (apiConfig *Config) PostPropertyTypesHandler(w http.ResponseWriter, r *http.Request, user User) {

	body := struct {
		Name string `json:"name"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error decoding request body. err: %v", err))
		return
	}
	if body.Name == "" {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Enter the property type name.")
		return
	}
	propertyType, err := apiConfig.DB.CreatePropertyType(r.Context(), body.Name)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			helpers.RespondWithError(w, http.StatusBadRequest, "property_type already exists")
			return
		}
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating property type. err: %v", err))
		return
	}
	log.Println(propertyType)
	helpers.RespondWithJson(w, http.StatusOK, DbPropertyTypeToModelsPropertyType(propertyType))

}

func (apiConfig *Config) GetPropertyTypesHandler(w http.ResponseWriter, r *http.Request) {
	propertyTypes, err := apiConfig.DB.GetPropertyTypes(r.Context())
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting property types. err: %v", err))
		return
	}
	helpers.RespondWithJson(w, http.StatusOK, DbPropertyTypesToModelsPropertyTypes(propertyTypes))

}

func (apiConfig *Config) GetPropertyTypeHandler(w http.ResponseWriter, r *http.Request) {

	name := chi.URLParam(r, "NAME")
	if name == "" {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Property type name is empty.")
		return
	}
	propertyType, err := apiConfig.DB.GetPropertyTypeWithName(r.Context(), name)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting property type. err: %v", err))
		return
	}

	helpers.RespondWithJson(w, http.StatusOK, DbPropertyTypeToModelsPropertyType(propertyType))

}
