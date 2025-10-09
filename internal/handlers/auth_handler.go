package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/muhammadolammi/rentradar/internal/auth"
	"github.com/muhammadolammi/rentradar/internal/database"
	"github.com/muhammadolammi/rentradar/internal/helpers"
	"golang.org/x/crypto/bcrypt"
)

func (apiConfig *Config) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Role        string `json:"role"`
		PhoneNumber string `json:"phone_number"`
		CompanyName string `json:"company_name"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding request body. err: %v", err))
		return
	}
	if body.Email == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Enter a mail.")
		return
	}
	if body.Password == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Enter a password.")
		return
	}

	// check if user exist
	userExist, err := apiConfig.DB.UserExists(r.Context(), body.Email)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error validating user. err: %v", err))
		return
	}
	if userExist {
		helpers.RespondWithError(w, http.StatusBadRequest, "User already exist. Login")
		return
	}
	// Validate role and role in enum
	if body.Role == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Enter the user role.")
		return
	}
	if body.Role != "user" && body.Role != "agent" && body.Role != "admin" && body.Role != "landlord" {
		helpers.RespondWithError(w, http.StatusBadRequest, "User  role must be one of (user, agent, landlord or admin)")
		return
	}

	if body.Role == "admin" {
		helpers.RespondWithError(w, http.StatusUnauthorized, "admin sign up not allowed")
		return
	}
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error hashing password. err: %v", err))
		return
	}
	// construct phone number is available
	userPhoneNumber := sql.NullString{String: "", Valid: false}
	if body.PhoneNumber != "" {
		userPhoneNumber = sql.NullString{String: body.PhoneNumber, Valid: true}
	}
	body.Email = strings.ToLower(strings.TrimSpace(body.Email))

	user, err := apiConfig.DB.CreateUser(r.Context(), database.CreateUserParams{
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		Email:       body.Email,
		Password:    string(hashedPassword),
		Role:        body.Role,
		PhoneNumber: userPhoneNumber,
	})
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating user. err: %v", err))
		return
	}
	// Update the company name if user role is agent
	if body.Role == "agent" {
		// company must exist
		if body.CompanyName == "" {
			helpers.RespondWithError(w, http.StatusBadRequest, "Enter the company name if registering as an agent")
			return
		}
		err = apiConfig.DB.UpdateUserCompanyName(r.Context(), database.UpdateUserCompanyNameParams{
			ID:          user.ID,
			CompanyName: sql.NullString{Valid: true, String: body.CompanyName},
		})
		if err != nil {
			helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating user. err: %v", err))
			return
		}
	}
	helpers.RespondWithJson(w, 200, "signup successful")
}

func (apiConfig *Config) LoginHandler(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding body from http request. err: %v", err))
		return
	}
	if body.Email == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Enter a mail.")
		return
	}
	if body.Password == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "Enter a password.")
		return
	}
	body.Email = strings.ToLower(strings.TrimSpace(body.Email))

	userExist, err := apiConfig.DB.UserExists(r.Context(), body.Email)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error validating user. err: %v", err))
		return
	}
	if !userExist {
		helpers.RespondWithError(w, http.StatusUnauthorized, "No User with this mail. Signup")
		return
	}

	user, err := apiConfig.DB.GetUserWithEmail(r.Context(), body.Email)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting user. err: %v", err))
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		if strings.Contains(err.Error(), `hashedPassword is not the hash of the given password`) {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Wrong password.")
			return
		}
		helpers.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf(" err: %v", err))
		return
	}
	// create refresh token

	err = auth.CreateRefreshToken([]byte(apiConfig.JWTKEY), user.ID, 24*7*6, w, apiConfig.DB)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating refresh token. err: %v", err))
		return
	}
	access_token, err := auth.MakeJwtTokenString([]byte(apiConfig.JWTKEY), user.ID.String(), "access_token", 15)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error making jwt token. err: %v", err))
		return
	}
	response := struct {
		AccessToken string `json:"access_token"`
		// ExpiresAt   time.Time `json:"expires_at"`
	}{
		AccessToken: access_token,
	}

	helpers.RespondWithJson(w, 200, response)
}

func (apiConfig *Config) PasswordChangeHandler(w http.ResponseWriter, r *http.Request) {

	body := struct {
		Email       string `json:"email"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding request body. err: %v", err))
		return
	}
	if body.Email == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Enter a mail. err: %v", err))
		return
	}
	if body.OldPassword == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Enter a password. err: %v", err))
		return
	}
	if body.NewPassword == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Enter a new password. err: %v", err))
		return
	}

	user, err := apiConfig.DB.GetUserWithEmail(r.Context(), body.Email)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting user. err: %v", err))
		return
	}
	// AUTHENTICATE THE USER
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.OldPassword))
	if err != nil {
		if strings.Contains(err.Error(), `hashedPassword is not the hash of the given password`) {
			helpers.RespondWithError(w, http.StatusUnauthorized, "Wrong password.")
			return
		}
		helpers.RespondWithError(w, http.StatusUnauthorized, fmt.Sprintf(" err: %v", err))
		return
	}
	// UPDATE THE PASSWORD
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 10)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error hashing password. err: %v", err))
		return
	}

	err = apiConfig.DB.UpdatePassword(r.Context(), database.UpdatePasswordParams{
		Email:    body.Email,
		Password: string(newHashedPassword),
	})
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error updating password. err: %v", err))
		return
	}
	err = auth.UpdateRefreshToken([]byte(apiConfig.JWTKEY), user.ID, 24*7*60, w, apiConfig.DB)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error updating refresh token. err: %v", err))
		return
	}
	helpers.RespondWithJson(w, 200, "Password Updated")
}

func (apiConfig *Config) GetUserHandler(w http.ResponseWriter, r *http.Request, user User) {

	helpers.RespondWithJson(w, 200, user)
}

func (apiConfig *Config) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	refreshtoken, err := r.Cookie("refresh_token")
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting refreshToken, Try login again. err: %v", err))
		return
	}

	refreshclaims := &jwt.RegisteredClaims{}

	_, err = jwt.ParseWithClaims(
		refreshtoken.Value,
		refreshclaims,
		func(token *jwt.Token) (interface{}, error) { return []byte(apiConfig.JWTKEY), nil },
	)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error parsing jwt claims. err: %v", err))
		return
	}
	// Make sure refresh token exist in db
	refreshTokenExists, err := apiConfig.DB.RefreshTokenExists(context.Background(), refreshtoken.Value)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error checking refresh token in db. err: %v", err))
		return
	}
	if !refreshTokenExists {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("refresh token doesnt exist. err: %v", err))
		return

	}
	userIdString := refreshclaims.Issuer
	userId, err := uuid.Parse(userIdString)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error parsing user id, err: %v", err))
		return
	}

	user, err := apiConfig.DB.GetUser(r.Context(), userId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error getting user with id, err: %v", err))
		return
	}

	refreshExpiration := refreshclaims.ExpiresAt.Time

	if refreshExpiration.Before(time.Now().UTC()) {
		helpers.RespondWithError(w, http.StatusUnauthorized, "refresh token expired")
		return
	}

	access_token, err := auth.MakeJwtTokenString([]byte(apiConfig.JWTKEY), user.ID.String(), "access_token", 15)
	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error making jwt token. err: %v", err))
		return
	}
	response := struct {
		AccessToken string `json:"access_token"`
		// ExpiresAt   time.Time `json:"expires_at"`
	}{
		AccessToken: access_token,
	}
	helpers.RespondWithJson(w, 200, response)

}

func (apiConfig *Config) Validate(w http.ResponseWriter, r *http.Request) {

	helpers.RespondWithJson(w, 200, "logged in")
}
