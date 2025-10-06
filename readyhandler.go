package main

import "net/http"

func successResponse(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, http.StatusOK, "server is ready for a success response")
}

func errorResponse(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "server is ready for a error response")
}
