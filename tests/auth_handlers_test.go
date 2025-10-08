package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterLoginAndRefresh(t *testing.T) {
	env := SetupTestEnv(t)
	t.Logf("-- Registering User ")

	// ---------- REGISTER ----------
	registerBody := map[string]string{
		"email":        "testuser@example.com",
		"password":     "StrongPass123",
		"first_name":   "John",
		"last_name":    "Doe",
		"role":         "user",
		"phone_number": "1234567890",
	}
	registerJSON, _ := json.Marshal(registerBody)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(registerJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-KEY", env.App.APIKEY)

	w := httptest.NewRecorder()

	// env.App.RegisterHandler(w, req)
	env.Router.ServeHTTP(w, req)

	// If the user already exists, continue
	if w.Code == http.StatusBadRequest && strings.Contains(w.Body.String(), "User already exist") {
		t.Log("User already exists — continuing test.")
	} else if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
	t.Logf("✅ Successfully registered ")
	t.Logf("--- Loging User In ")

	// ---------- LOGIN ----------
	loginBody := map[string]string{
		"email":    "testuser@example.com",
		"password": "StrongPass123",
	}
	loginJSON, _ := json.Marshal(loginBody)

	req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-KEY", env.App.APIKEY)

	w = httptest.NewRecorder()

	// env.App.LoginHandler(w, req)
	env.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var loginResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("error parsing login response: %v", err)
	}
	if loginResp.AccessToken == "" {
		t.Fatal("access_token not found in response")
	}
	t.Logf("✅ Successfully Logged In")

	t.Logf("--- Refreshing token")

	// ---------- REFRESH ----------
	// get refresh cookie set by LoginHandler
	cookies := w.Result().Cookies()
	var refreshCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "refresh_token" {
			refreshCookie = c
			break
		}
	}

	if refreshCookie == nil {
		t.Fatal("refresh_token cookie not found after login")
	}

	req = httptest.NewRequest(http.MethodPost, "/refresh", nil)
	req.Header.Set("API-KEY", env.App.APIKEY)

	req.AddCookie(refreshCookie)
	w = httptest.NewRecorder()

	env.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 from refresh, got %d, body: %s", w.Code, w.Body.String())
	}
	t.Logf("✅ Successfully refreshed tokens.")

}
