package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/muhammadolammi/rentradar/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (apiConfig *Config) registerHandler(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Role        string `json:"role"`
		PhoneNumber string `json:"phone_number"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding body from http request. err: %v", err))
		return
	}
	if body.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Enter a mail.")
		return
	}
	if body.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Enter a password.")
		return
	}

	// check if user exist
	userExist, err := apiConfig.DB.UserExists(r.Context(), body.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error validating user. err: %v", err))
		return
	}
	if userExist {
		respondWithError(w, http.StatusBadRequest, "User already exist. Login")
		return
	}
	// Validate role and role in enum
	if body.Role == "" {
		respondWithError(w, http.StatusBadRequest, "Enter the user role.")
		return
	}
	if body.Role != "user" && body.Role != "agent" && body.Role != "admin" && body.Role != "landlord" {
		respondWithError(w, http.StatusBadRequest, "User  role must be one of (user, agent, landlord or admin)")
		return
	}
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error hashing password. err: %v", err))
		return
	}
	// construct phone number is available
	userPhoneNumber := sql.NullString{String: "", Valid: false}
	if body.Password != "" {
		userPhoneNumber = sql.NullString{String: body.Password, Valid: true}

	}

	_, err = apiConfig.DB.CreateUser(r.Context(), database.CreateUserParams{
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		Email:       body.Email,
		Password:    string(hashedPassword),
		Role:        body.Role,
		PhoneNumber: userPhoneNumber,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating user. err: %v", err))
		return
	}

	respondWithJson(w, 200, "signup successful")
}

func (apiConfig *Config) loginHandler(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding body from http request. err: %v", err))
		return
	}
	if body.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Enter a mail.")
		return
	}
	if body.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Enter a password.")
		return
	}

	userExist, err := apiConfig.DB.UserExists(r.Context(), body.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error validating user. err: %v", err))
		return
	}
	if !userExist {
		respondWithError(w, http.StatusUnauthorized, "No User with this mail. Signup")
		return
	}

	user, err := apiConfig.DB.GetUserWithEmail(r.Context(), body.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting user. err: %v", err))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		if strings.Contains(err.Error(), `hashedPassword is not the hash of the given password`) {
			respondWithError(w, http.StatusUnauthorized, "Wrong password.")
			return
		}
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf(" err: %v", err))
		return
	}

	respondWithJson(w, 200, "logged in")
}
