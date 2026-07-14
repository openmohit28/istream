package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"istream/backend/internal/config"
	"istream/backend/internal/database"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		url = "postgres://mohitrawat@localhost:5432/istream_test?sslmode=disable"
	}
	var err error
	testDB, err = database.Connect(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "test database unavailable: %v\n", err)
		os.Exit(1)
	}
	if err := database.Migrate(testDB); err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
		os.Exit(1)
	}
	code := m.Run()
	testDB.Close()
	os.Exit(code)
}

func newTestServer(t *testing.T) http.Handler {
	t.Helper()
	if _, err := testDB.Exec(`TRUNCATE users CASCADE`); err != nil {
		t.Fatalf("truncate users: %v", err)
	}
	cfg := config.Config{
		JWTSecret:      "test-secret",
		FrontendOrigin: "http://localhost:5173",
	}
	return NewServer(testDB, cfg)
}

func doJSON(t *testing.T, srv http.Handler, method, path string, body any, token string) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	return rec
}

func decodeAuth(t *testing.T, rec *httptest.ResponseRecorder) (token string, user map[string]any) {
	t.Helper()
	var resp struct {
		Token string         `json:"token"`
		User  map[string]any `json:"user"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode auth response: %v (body: %s)", err, rec.Body.String())
	}
	return resp.Token, resp.User
}

var validRegister = map[string]string{
	"email":    "mohit@example.com",
	"name":     "Mohit",
	"password": "supersecret1",
}

func TestRegisterSuccess(t *testing.T) {
	srv := newTestServer(t)
	rec := doJSON(t, srv, "POST", "/api/auth/register", validRegister, "")
	if rec.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", rec.Code, rec.Body.String())
	}
	token, user := decodeAuth(t, rec)
	if token == "" {
		t.Error("expected a token")
	}
	if user["email"] != "mohit@example.com" {
		t.Errorf("want email mohit@example.com, got %v", user["email"])
	}
	if _, leaked := user["passwordHash"]; leaked {
		t.Error("password hash must not appear in responses")
	}
	if _, leaked := user["password_hash"]; leaked {
		t.Error("password hash must not appear in responses")
	}
}

func TestRegisterNormalizesEmail(t *testing.T) {
	srv := newTestServer(t)
	body := map[string]string{"email": "  Mohit@Example.COM ", "name": "Mohit", "password": "supersecret1"}
	rec := doJSON(t, srv, "POST", "/api/auth/register", body, "")
	if rec.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", rec.Code, rec.Body.String())
	}
	_, user := decodeAuth(t, rec)
	if user["email"] != "mohit@example.com" {
		t.Errorf("email not normalized: %v", user["email"])
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	srv := newTestServer(t)
	doJSON(t, srv, "POST", "/api/auth/register", validRegister, "")
	rec := doJSON(t, srv, "POST", "/api/auth/register", validRegister, "")
	if rec.Code != http.StatusConflict {
		t.Fatalf("want 409, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestRegisterValidation(t *testing.T) {
	srv := newTestServer(t)
	cases := []struct {
		name string
		body map[string]string
	}{
		{"invalid email", map[string]string{"email": "not-an-email", "name": "M", "password": "supersecret1"}},
		{"empty name", map[string]string{"email": "a@b.com", "name": "  ", "password": "supersecret1"}},
		{"short password", map[string]string{"email": "a@b.com", "name": "M", "password": "short"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := doJSON(t, srv, "POST", "/api/auth/register", tc.body, "")
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("want 400, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestLoginSuccess(t *testing.T) {
	srv := newTestServer(t)
	doJSON(t, srv, "POST", "/api/auth/register", validRegister, "")
	rec := doJSON(t, srv, "POST", "/api/auth/login",
		map[string]string{"email": "mohit@example.com", "password": "supersecret1"}, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", rec.Code, rec.Body.String())
	}
	token, _ := decodeAuth(t, rec)
	if token == "" {
		t.Error("expected a token")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	srv := newTestServer(t)
	doJSON(t, srv, "POST", "/api/auth/register", validRegister, "")
	rec := doJSON(t, srv, "POST", "/api/auth/login",
		map[string]string{"email": "mohit@example.com", "password": "wrongpassword"}, "")
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestLoginUnknownEmail(t *testing.T) {
	srv := newTestServer(t)
	rec := doJSON(t, srv, "POST", "/api/auth/login",
		map[string]string{"email": "ghost@example.com", "password": "supersecret1"}, "")
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestMeWithValidToken(t *testing.T) {
	srv := newTestServer(t)
	reg := doJSON(t, srv, "POST", "/api/auth/register", validRegister, "")
	token, _ := decodeAuth(t, reg)

	rec := doJSON(t, srv, "GET", "/api/auth/me", nil, token)
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var user map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &user); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if user["email"] != "mohit@example.com" {
		t.Errorf("want email mohit@example.com, got %v", user["email"])
	}
}

func TestMeWithoutToken(t *testing.T) {
	srv := newTestServer(t)
	rec := doJSON(t, srv, "GET", "/api/auth/me", nil, "")
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestMeWithGarbageToken(t *testing.T) {
	srv := newTestServer(t)
	rec := doJSON(t, srv, "GET", "/api/auth/me", nil, "not.a.token")
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("want 401, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHealth(t *testing.T) {
	srv := newTestServer(t)
	rec := doJSON(t, srv, "GET", "/api/health", nil, "")
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
}
