package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestAlertEndpoints tests creating and retrieving alerts for a user.
func TestAlertEndpoints(t *testing.T) {
	env := SetupTestEnv(t)

	// ---------- Register a user ----------
	t.Log("--- Registering user")
	registerBody := map[string]string{
		"email":        "alertuser@example.com",
		"password":     "StrongPass123",
		"first_name":   "Alert",
		"last_name":    "Tester",
		"role":         "user",
		"phone_number": "08000000000",
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
		"email":    "alertuser@example.com",
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

	// ---------- Create Alert ----------
	t.Log("--- Creating Alert")
	alertBody := map[string]any{
		"min_price":      100000,
		"max_price":      300000,
		"location":       "Lagos",
		"property_type":  "apartment",
		"contact_method": "email",
	}
	alertJSON, _ := json.Marshal(alertBody)
	req = httptest.NewRequest(http.MethodPost, "/alerts", bytes.NewBuffer(alertJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+loginResp.AccessToken)
	req.Header.Set("API-KEY", env.App.APIKEY)

	w = httptest.NewRecorder()
	env.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 from PostAlertHandler, got %d, body: %s", w.Code, w.Body.String())
	}

	var alertResp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &alertResp); err != nil {
		t.Fatalf("error parsing create alert response: %v", err)
	}
	t.Log("✅ Successfully Created Alert")

	// ---------- Get User Alerts ----------
	t.Log("--- Getting User Alerts")
	req = httptest.NewRequest(http.MethodGet, "/alerts", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.AccessToken)
	req.Header.Set("API-KEY", env.App.APIKEY)

	w = httptest.NewRecorder()
	env.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 from GetUserAlertsHandler, got %d, body: %s", w.Code, w.Body.String())
	}

	var alertsResp []map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &alertsResp); err != nil {
		t.Fatalf("error parsing get user alerts response: %v", err)
	}

	if len(alertsResp) == 0 {
		t.Fatalf("expected at least 1 alert, got 0")
	}
	t.Logf("✅ Successfully retrieved %d alert(s)", len(alertsResp))
}
