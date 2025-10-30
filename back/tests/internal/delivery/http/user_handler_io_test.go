package httptestpkg

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUsers_IO(t *testing.T) {
	r := NewRouter()

	// empty body on create
	rr := doReq(t, r, http.MethodPost, "/api/users", nil)
	if rr.Code == http.StatusOK || rr.Code == http.StatusCreated {
		t.Fatalf("expected validation error on empty body, got %d", rr.Code)
	}

	// invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code == http.StatusOK || rr.Code == http.StatusCreated {
		t.Fatalf("expected error on invalid json, got %d", rr.Code)
	}

	// create valid
	payload := map[string]any{
		"name":  "Alice",
		"email": "alice@example.com",
	}
	rr = doReq(t, r, http.MethodPost, "/api/users", payload)
	if rr.Code != http.StatusCreated && rr.Code != http.StatusOK {
		t.Fatalf("expected create success, got %d: %s", rr.Code, rr.Body.String())
	}

	// get list
	rr = doReq(t, r, http.MethodGet, "/api/users", nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected list 200, got %d", rr.Code)
	}

	// get by bad id
	rr = doReq(t, r, http.MethodGet, "/api/users/bad-id", nil)
	if rr.Code == http.StatusOK {
		t.Fatalf("expected error on bad id, got %d", rr.Code)
	}

	// update with empty body
	rr = doReq(t, r, http.MethodPut, "/api/users/1", nil)
	if rr.Code == http.StatusOK {
		t.Fatalf("expected validation error on empty update, got %d", rr.Code)
	}

	// update valid
	upd := map[string]any{"name": "Alice 2"}
	rr = doReq(t, r, http.MethodPut, "/api/users/1", upd)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected update 200, got %d: %s", rr.Code, rr.Body.String())
	}

	// delete invalid id
	rr = doReq(t, r, http.MethodDelete, "/api/users/NaN", nil)
	if rr.Code == http.StatusOK || rr.Code == http.StatusNoContent {
		t.Fatalf("expected error on delete invalid id, got %d", rr.Code)
	}

	// delete ok
	rr = doReq(t, r, http.MethodDelete, "/api/users/1", nil)
	if rr.Code != http.StatusOK && rr.Code != http.StatusNoContent {
		t.Fatalf("expected delete success, got %d: %s", rr.Code, rr.Body.String())
	}
}
