package httptestpkg

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type creds struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshReq struct {
	RefreshToken string `json:"refresh_token"`
}

func TestUser_Register_Login_Refresh_And_Me_CRUD(t *testing.T) {
	// Arrange: spin up router with minimal deps via NewRouter (black-box)
	r := NewRouter()

	// 1) Register valid
	reg := creds{Email: "user@example.com", Password: "Str0ngPass!"}
	b, err := json.Marshal(reg)
	if err != nil {
		t.Fatalf("failed to marshal register request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	if rw.Code != http.StatusCreated && rw.Code != http.StatusOK {
		t.Fatalf("register expected 201/200, got %d, body=%s", rw.Code, rw.Body.String())
	}

	// 2) Register invalid (weak password)
	regWeak := creds{Email: "weak@example.com", Password: "123"}
	b, err = json.Marshal(regWeak)
	if err != nil {
		t.Fatalf("failed to marshal weak password request: %v", err)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/v1/user/register", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	if rw.Code == http.StatusOK || rw.Code == http.StatusCreated {
		t.Fatalf("expected validation error on weak password, got %d", rw.Code)
	}

	// 3) Login valid
	b, err = json.Marshal(reg)
	if err != nil {
		t.Fatalf("failed to marshal login request: %v", err)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/v1/user/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("login expected 200, got %d, body=%s", rw.Code, rw.Body.String())
	}

	var loginResp map[string]any
	if err := json.Unmarshal(rw.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("failed to unmarshal login response: %v", err)
	}

	access, _ := loginResp["access_token"].(string)
	refresh, _ := loginResp["refresh_token"].(string)
	if access == "" || refresh == "" {
		t.Fatalf("expected access and refresh tokens in login response")
	}

	// 4) Me GET with access
	req = httptest.NewRequest(http.MethodGet, "/api/v1/user/me", nil)
	req.Header.Set("Authorization", "Bearer "+access)
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("me GET expected 200, got %d", rw.Code)
	}

	// 5) Me PUT update
	update := map[string]any{"name": "Updated User"}
	b, err = json.Marshal(update)
	if err != nil {
		t.Fatalf("failed to marshal update request: %v", err)
	}

	req = httptest.NewRequest(http.MethodPut, "/api/v1/user/me", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+access)
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("me PUT expected 200, got %d", rw.Code)
	}

	// 6) Refresh token
	b, err = json.Marshal(refreshReq{RefreshToken: refresh})
	if err != nil {
		t.Fatalf("failed to marshal refresh request: %v", err)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/v1/user/refresh", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("refresh expected 200, got %d", rw.Code)
	}

	// 7) Me DELETE account
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/user/me", nil)
	req.Header.Set("Authorization", "Bearer "+access)
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK && rw.Code != http.StatusNoContent {
		t.Fatalf("me DELETE expected 200/204, got %d", rw.Code)
	}
}
