package main

import "net/http"

// Middleware to check for the API key in the authorization header for all requests.
func (apiConfig *Config) verifyApiKey() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			api_key := r.Header.Get("Authorization")
			if api_key != apiConfig.APIKEY {
				respondWithError(w, http.StatusUnauthorized, "Invalid Api key")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
