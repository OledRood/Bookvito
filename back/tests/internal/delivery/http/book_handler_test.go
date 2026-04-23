package httptestpkg

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// NOTE: These are I/O focused tests based on expected routes and payloads.
// They avoid relying on internal implementation details.

func TestHealthEndpoint_OK(t *testing.T) {
	// Arrange: use the test router stub
	r := NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// Act: send through the router
	r.ServeHTTP(w, req)

	// Assert: expect 200 OK from the health endpoint
	if status := w.Code; status != http.StatusOK {
		t.Fatalf("health endpoint expected 200, got %d", status)
	}

	// Verify response contains expected status
	body := w.Body.String()
	if body != `{"status":"ok"}` {
		t.Errorf("unexpected response body: %s", body)
	}
}
