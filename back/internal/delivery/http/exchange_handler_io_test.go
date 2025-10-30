package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func doReq(t *testing.T, h http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf *bytes.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		buf = bytes.NewReader(b)
	} else {
		buf = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, buf)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}

func TestExchange_IO(t *testing.T) {
	r := NewRouter()

	// create: empty body
	rr := doReq(t, r, http.MethodPost, "/api/exchanges", nil)
	if rr.Code == http.StatusOK || rr.Code == http.StatusCreated {
		t.Fatalf("expected validation error on empty body, got %d", rr.Code)
	}

	// create: invalid json
	req := httptest.NewRequest(http.MethodPost, "/api/exchanges", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code == http.StatusOK || rr.Code == http.StatusCreated {
		t.Fatalf("expected error on invalid json, got %d", rr.Code)
	}

	// create: valid
	payload := map[string]any{
		"book_id": 1,
		"user_id": 1,
	}
	rr = doReq(t, r, http.MethodPost, "/api/exchanges", payload)
	if rr.Code != http.StatusCreated && rr.Code != http.StatusOK {
		t.Fatalf("expected create success, got %d: %s", rr.Code, rr.Body.String())
	}

	// list
	rr = doReq(t, r, http.MethodGet, "/api/exchanges", nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected list 200, got %d", rr.Code)
	}

	// get by invalid id
	rr = doReq(t, r, http.MethodGet, "/api/exchanges/bad", nil)
	if rr.Code == http.StatusOK {
		t.Fatalf("expected error on bad id, got %d", rr.Code)
	}

	// update: empty body
	rr = doReq(t, r, http.MethodPut, "/api/exchanges/1", nil)
	if rr.Code == http.StatusOK {
		t.Fatalf("expected validation error, got %d", rr.Code)
	}

	// update: valid
	upd := map[string]any{"status": "closed"}
	rr = doReq(t, r, http.MethodPut, "/api/exchanges/1", upd)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected update 200, got %d: %s", rr.Code, rr.Body.String())
	}

	// delete: invalid id
	rr = doReq(t, r, http.MethodDelete, "/api/exchanges/NaN", nil)
	if rr.Code == http.StatusOK || rr.Code == http.StatusNoContent {
		t.Fatalf("expected error on delete invalid id, got %d", rr.Code)
	}

	// delete ok
	rr = doReq(t, r, http.MethodDelete, "/api/exchanges/1", nil)
	if rr.Code != http.StatusOK && rr.Code != http.StatusNoContent {
		t.Fatalf("expected delete success, got %d: %s", rr.Code, rr.Body.String())
	}
}
