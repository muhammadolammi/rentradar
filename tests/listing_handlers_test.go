package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateAndGetListings(t *testing.T) {
	env := SetupTestEnv(t)
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
	w := httptest.NewRecorder()

	env.App.RegisterHandler(w, req)

	// If the user already exists, continue
	if w.Code == http.StatusBadRequest && strings.Contains(w.Body.String(), "User already exist") {
		t.Log("User already exists — continuing test.")
	} else if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	// ---------- LOGIN AS AGENT ----------
	loginBody := map[string]string{
		"email":    "agent@example.com",
		"password": "StrongPass123",
	}
	loginJSON, _ := json.Marshal(loginBody)

	req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	env.App.LoginHandler(w, req)

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

	// ---------- CREATE LISTING ----------
	postBody := map[string]any{
		"description": "A beautiful 3-bedroom apartment with modern amenities.",
		"title":       "Modern Apartment",
		"rent_type":   "monthly",
		"house_type":  "apartment",
		"images":      []string{"img1.jpg", "img2.jpg"},
		"price":       250000,
		"location":    "Lagos",
	}

	postJSON, _ := json.Marshal(postBody)
	req = httptest.NewRequest(http.MethodPost, "/listings", bytes.NewBuffer(postJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+loginResp.AccessToken)
	w = httptest.NewRecorder()

	handler := env.App.AuthMiddleware([]byte(env.App.JWTKEY), env.App.PostListingsHandler)
	handler(w, req)

	if w.Code == http.StatusBadRequest && strings.Contains(w.Body.String(), "listing already exists") {
		t.Log("Listing already exists — continuing test.")
	} else if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var postResp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &postResp); err != nil {
		t.Fatalf("error parsing create listing response: %v", err)
	}

	if _, ok := postResp["title"]; !ok {
		t.Fatal("expected title field in create listing response")
	}

	// ---------- GET LISTINGS ----------
	getBody := map[string]any{
		"page":  1,
		"limit": 5,
	}
	getJSON, _ := json.Marshal(getBody)
	req = httptest.NewRequest(http.MethodPost, "/listings?location=Lagos&type=apartment", bytes.NewBuffer(getJSON))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	env.App.GetListingsHandler(w, req)

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
}
