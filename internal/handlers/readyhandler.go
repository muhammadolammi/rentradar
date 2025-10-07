package handlers

import (
	"net/http"

	"github.com/muhammadolammi/rentradar/internal/helpers"
)

func SuccessResponse(w http.ResponseWriter, r *http.Request) {
	helpers.RespondWithJson(w, http.StatusOK, "server is ready for a success response")
}

func ErrorResponse(w http.ResponseWriter, r *http.Request) {
	helpers.RespondWithError(w, http.StatusInternalServerError, "server is ready for a error response")
}
