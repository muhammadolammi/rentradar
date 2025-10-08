package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestListingsEndpoints(t *testing.T) {
	env := SetupTestEnv(t)
	t.Log("--- Registering agent")

	// ---------- REGISTER ----------
	registerBody := map[string]string{
		"email":        "agent@example.com",
		"password":     "StrongPass123",
		"first_name":   "Akek",
		"last_name":    "kek",
		"role":         "agent",
		"phone_number": "12345678970",
		"company_name": "rent_radar",
	}
	registerJSON, _ := json.Marshal(registerBody)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(registerJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-KEY", env.App.APIKEY)

	w := httptest.NewRecorder()

	env.Router.ServeHTTP(w, req) // ✅ goes through router + middlewares

	if w.Code == http.StatusBadRequest && strings.Contains(w.Body.String(), "User already exist") {
		t.Log("User already exists — continuing test.")
	} else if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
	t.Log("✅ Successfully Registered ")

	t.Log("--- Logging Agent in.")

	// ---------- LOGIN AS AGENT ----------
	loginBody := map[string]string{
		"email":    "agent@example.com",
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
	t.Log("✅ Successfully Logged In ")

	// ---------- POST PROPERTY TYPE ----------
	t.Log("-- Posting Property Type")

	postBody := map[string]any{
		"name": "apartment",
	}
	postJSON, _ := json.Marshal(postBody)

	req = httptest.NewRequest(http.MethodPost, "/property_types", bytes.NewBuffer(postJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-KEY", env.App.APIKEY)

	req.Header.Set("Authorization", "Bearer "+loginResp.AccessToken)
	req.Header.Set("SUDO-KEY", env.App.SUDOKEY) // ✅ pass sudo key
	w = httptest.NewRecorder()

	env.Router.ServeHTTP(w, req)

	var postResp map[string]any

	if w.Code == http.StatusBadRequest && strings.Contains(w.Body.String(), "property_type already exists") {
		t.Log("PropertyType already exists — continuing test.")
	} else if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
	// else {
	// 	if err := json.Unmarshal(w.Body.Bytes(), &postResp); err != nil {
	// 		t.Fatalf("error parsing create property type response: %v", err)
	// 	}
	// 	if _, ok := postResp["name"]; !ok {
	// 		t.Fatal("expected name field in create property type response")
	// 	}
	// 	propertyTypeId = postResp["id"]
	// 	t.Log("✅ Successfully Created Property Type")
	// }
	// ---------- POST PROPERTY TYPE ----------
	// todo
	t.Log("-- Getting Property Type")
	req = httptest.NewRequest(http.MethodGet, "/property_types/apartment", bytes.NewBuffer(postJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-KEY", env.App.APIKEY)
	req.Header.Set("Authorization", "Bearer "+loginResp.AccessToken)
	w = httptest.NewRecorder()
	env.Router.ServeHTTP(w, req)
	getresponse := map[string]any{}
	if err := json.Unmarshal(w.Body.Bytes(), &getresponse); err != nil {
		t.Fatalf("error parsing get property type response: %v", err)
	}
	if _, ok := getresponse["name"]; !ok {
		t.Fatal("expected name field in get property type response")
	}
	if _, ok := getresponse["id"]; !ok {
		t.Fatal("expected id field in get property type response")
	}
	propertyTypeId, ok := getresponse["id"].(string)
	if !ok || propertyTypeId == "" {
		t.Fatalf("expected valid propertyTypeId string, got %#v", getresponse["id"])
	}
	t.Logf("✅ Successfully Retrieved Property Type: %s", propertyTypeId)
	t.Log("✅ Successfully Retrieved Property Type")

	// ---------- CREATE LISTING ----------
	t.Log("--- Creating Listings")

	postBody = map[string]any{
		"description":      "A beautiful 3-bedroom apartment with modern amenities.",
		"title":            "Modern Apartment",
		"rent_type":        "monthly",
		"property_type_id": propertyTypeId,
		"images":           []string{"img1.jpg", "img2.jpg"},
		"price":            250000,
		"location":         "Lagos",
	}
	postJSON, _ = json.Marshal(postBody)

	req = httptest.NewRequest(http.MethodPost, "/listings", bytes.NewBuffer(postJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-KEY", env.App.APIKEY)

	req.Header.Set("Authorization", "Bearer "+loginResp.AccessToken)
	w = httptest.NewRecorder()

	env.Router.ServeHTTP(w, req)

	if w.Code == http.StatusBadRequest && strings.Contains(w.Body.String(), "listing already exists") {
		t.Log("Listing already exists — continuing test.")
	} else if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	postResp = map[string]any{}
	if err := json.Unmarshal(w.Body.Bytes(), &postResp); err != nil {
		t.Fatalf("error parsing create listing response: %v", err)
	}

	if _, ok := postResp["title"]; !ok {
		t.Fatal("expected title field in create listing response")
	}
	t.Log("✅ Successfully Created Listing ")

	// ---------- GET LISTINGS ----------
	t.Log("--- Getting Listings")

	req = httptest.NewRequest(http.MethodGet, "/listings?location=Lagos&property_type_name=apartment", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-KEY", env.App.APIKEY)

	w = httptest.NewRecorder()

	env.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 from GetListings, got %d, body: %s", w.Code, w.Body.String())
	}

	var getResp []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &getResp); err != nil {
		t.Fatalf("error parsing get listings response: %v", err)
	}

	if len(getResp) == 0 {
		t.Fatalf("expected at least 1 listing in response")
	}

	t.Logf("✅ Successfully retrieved %d listing(s)", len(getResp))

	// ---------- GET SINGLE LISTING ----------
	t.Log("--- Getting Single Listing")

	listingID := getResp[0]["id"].(string)

	req = httptest.NewRequest(http.MethodGet, "/listings/"+listingID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("ID", listingID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	req.Header.Set("API-KEY", env.App.APIKEY)

	w = httptest.NewRecorder()

	env.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 from GetListing, got %d, body: %s", w.Code, w.Body.String())
	}

	var getListingResp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &getListingResp); err != nil {
		t.Fatalf("error parsing get listing response: %v", err)
	}
	if getListingResp["id"] == "" {
		t.Fatalf("expected id field in get listing response")
	}
	t.Log("✅ Successfully retrieved single listing")
}
