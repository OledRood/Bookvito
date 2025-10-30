package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type locationCreate struct {
	Name string `json:"name"`
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
}

func TestLocation_Create_Get_List_Validation(t *testing.T) {
	r := NewRouter()

	// Create valid
	b, _ := json.Marshal(locationCreate{Name: "Main Library", Lat: 40.1, Lng: -73.9})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/locations", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	if rw.Code != http.StatusCreated && rw.Code != http.StatusOK {
		t.Fatalf("create expected 201/200, got %d, body=%s", rw.Code, rw.Body.String())
	}

	// Create invalid (missing name)
	b, _ = json.Marshal(locationCreate{Name: "", Lat: 0, Lng: 0})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/locations", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	if rw.Code == http.StatusOK || rw.Code == http.StatusCreated {
		t.Fatalf("expected validation error for empty name, got %d", rw.Code)
	}

	// Get all
	req = httptest.NewRequest(http.MethodGet, "/api/v1/locations", nil)
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("list expected 200, got %d", rw.Code)
	}

	// Get by id
	req = httptest.NewRequest(http.MethodGet, "/api/v1/locations/1", nil)
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK && rw.Code != http.StatusNotFound {
		t.Fatalf("get by id expected 200 or 404, got %d", rw.Code)
	}
}
