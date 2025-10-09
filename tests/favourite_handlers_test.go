package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestFavoritesEndpoints tests saving a listing as favorite and fetching it.
func TestFavoritesEndpoints(t *testing.T) {
	env := SetupTestEnv(t)

	// ---------- Register a user ----------
	t.Log("--- Registering user")
	registerBody := map[string]string{
		"email":        "favuser@example.com",
		"password":     "StrongPass123",
		"first_name":   "Fav",
		"last_name":    "Tester",
		"role":         "agent",
		"company_name": "akek",
		"phone_number": "08000000001",
	}
	registerJSON, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(registerJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-KEY", env.App.APIKEY)

	w := httptest.NewRecorder()
	env.Router.ServeHTTP(w, req)

	if w.Code == http.StatusBadRequest && strings.Contains(w.Body.String(), "User already exist") {
		t.Log("User already exists — continuing test.")
	} else if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
	t.Log("✅ Successfully Registered")

	// ---------- Login ----------
	t.Log("--- Logging in user")
	loginBody := map[string]string{
		"email":    "favuser@example.com",
		"password": "StrongPass123",
	}
	loginJSON, _ := json.Marshal(loginBody)
	req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-KEY", env.App.APIKEY)

	w = httptest.NewRecorder()
	env.Router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Login failed: expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var loginResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("Error parsing login response: %v", err)
	}
	if loginResp.AccessToken == "" {
		t.Fatal("access_token missing in login response")
	}
	t.Log("✅ Successfully Logged In")

	// ---------- Create Listing ----------
	t.Log("--- Creating listing for favorite")
	listingBody := map[string]any{
		"title":         "Modern Duplex Apartment",
		"description":   "A clean duplex apartment in Lekki",
		"price":         500000,
		"location":      "Lekki",
		"property_type": "apartment",
		"images":        []string{"img1.jpg", "img2.jpg"},
	}
	listingJSON, _ := json.Marshal(listingBody)
	req = httptest.NewRequest(http.MethodPost, "/listings", bytes.NewBuffer(listingJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+loginResp.AccessToken)
	req.Header.Set("API-KEY", env.App.APIKEY)

	w = httptest.NewRecorder()
	env.Router.ServeHTTP(w, req)

	var listingResp map[string]any
	if w.Code == http.StatusBadRequest && strings.Contains(w.Body.String(), "listing already exists") {
		t.Log("Listing already exists — continuing test.")
		req = httptest.NewRequest(http.MethodGet, "/listings?location=Lekki", nil)
		req.Header.Set("Authorization", "Bearer "+loginResp.AccessToken)
		req.Header.Set("API-KEY", env.App.APIKEY)
		w = httptest.NewRecorder()
		env.Router.ServeHTTP(w, req)
		if err := json.Unmarshal(w.Body.Bytes(), &listingResp); err != nil {
			t.Fatalf("error parsing get listing: %v", err)
		}
	} else if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	} else {
		if err := json.Unmarshal(w.Body.Bytes(), &listingResp); err != nil {
			t.Fatalf("error parsing create listing response: %v", err)
		}
	}
	listingID := listingResp["id"].(string)
	t.Log("✅ Listing created")

	// ---------- Create Favorite ----------
	t.Log("--- Saving listing as favorite")
	favBody := map[string]any{
		"listing_id": listingID,
	}
	favJSON, _ := json.Marshal(favBody)
	req = httptest.NewRequest(http.MethodPost, "/favorites", bytes.NewBuffer(favJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+loginResp.AccessToken)
	req.Header.Set("API-KEY", env.App.APIKEY)

	w = httptest.NewRecorder()
	env.Router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 from PostFavoritesHandler, got %d, body: %s", w.Code, w.Body.String())
	}

	var favResp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &favResp); err != nil {
		t.Fatalf("error parsing create favorite response: %v", err)
	}
	t.Log("✅ Successfully saved listing as favorite")

	// ---------- Get Favorites ----------
	t.Log("--- Fetching user favorites")
	req = httptest.NewRequest(http.MethodGet, "/favorites", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.AccessToken)
	req.Header.Set("API-KEY", env.App.APIKEY)

	w = httptest.NewRecorder()
	env.Router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 from GetFavoritesHandler, got %d, body: %s", w.Code, w.Body.String())
	}

	var favsResp []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &favsResp); err != nil {
		t.Fatalf("error parsing get favorites response: %v", err)
	}

	if len(favsResp) == 0 {
		t.Fatalf("expected at least 1 favorite, got 0")
	}
	t.Logf("✅ Successfully retrieved %d favorite(s)", len(favsResp))
}
