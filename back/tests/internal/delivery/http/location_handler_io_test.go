package httptestpkg

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLocation_IO(t *testing.T) {
	r := NewRouter()

	// create: empty body
	rr := doReq(t, r, http.MethodPost, "/api/locations", nil)
	if rr.Code == http.StatusOK || rr.Code == http.StatusCreated {
		t.Fatalf("expected validation error on empty body, got %d", rr.Code)
	}

	// create: invalid json
	req := httptest.NewRequest(http.MethodPost, "/api/locations", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code == http.StatusOK || rr.Code == http.StatusCreated {
		t.Fatalf("expected error on invalid json, got %d", rr.Code)
	}

	// create: valid
	payload := map[string]any{
		"name": "Moscow",
	}
	rr = doReq(t, r, http.MethodPost, "/api/locations", payload)
	if rr.Code != http.StatusCreated && rr.Code != http.StatusOK {
		t.Fatalf("expected create success, got %d: %s", rr.Code, rr.Body.String())
	}

	// list
	rr = doReq(t, r, http.MethodGet, "/api/locations", nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected list 200, got %d", rr.Code)
	}

	// get by invalid id
	rr = doReq(t, r, http.MethodGet, "/api/locations/bad", nil)
	if rr.Code == http.StatusOK {
		t.Fatalf("expected error on bad id, got %d", rr.Code)
	}

	// update: empty body
	rr = doReq(t, r, http.MethodPut, "/api/locations/1", nil)
	if rr.Code == http.StatusOK {
		t.Fatalf("expected validation error, got %d", rr.Code)
	}

	// update: valid
	upd := map[string]any{"name": "Moscow City"}
	rr = doReq(t, r, http.MethodPut, "/api/locations/1", upd)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected update 200, got %d: %s", rr.Code, rr.Body.String())
	}

	// delete: invalid id
	rr = doReq(t, r, http.MethodDelete, "/api/locations/NaN", nil)
	if rr.Code == http.StatusOK || rr.Code == http.StatusNoContent {
		t.Fatalf("expected error on delete invalid id, got %d", rr.Code)
	}

	// delete ok
	rr = doReq(t, r, http.MethodDelete, "/api/locations/1", nil)
	if rr.Code != http.StatusOK && rr.Code != http.StatusNoContent {
		t.Fatalf("expected delete success, got %d: %s", rr.Code, rr.Body.String())
	}
}
