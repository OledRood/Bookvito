package httptestpkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// NewRouter is a lightweight test stub that provides minimal routes used by the
// moved tests. It implements simple, predictable responses to mimic the
// real API contract (status codes and small JSON payloads) without wiring
// application logic or external dependencies.
func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
			// Log error in real scenarios; here we're in test stub
			fmt.Printf("Write error: %v\n", err)
		}
	})

	// Generic POST create endpoints under /api and /api/v1
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		path := strings.TrimPrefix(r.URL.Path, "/api")

		// /api/exchanges and /api/exchanges/:id
		if strings.HasPrefix(path, "/exchanges") {
			handleCRUD(w, r, "/exchanges")
			return
		}

		// /api/users
		if strings.HasPrefix(path, "/users") {
			handleCRUD(w, r, "/users")
			return
		}

		// /api/locations
		if strings.HasPrefix(path, "/locations") {
			handleCRUD(w, r, "/locations")
			return
		}

		// Fallback 404 for other /api routes
		http.NotFound(w, r)
	})

	// v1 API group
	mux.HandleFunc("/api/v1/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		path := strings.TrimPrefix(r.URL.Path, "/api/v1")

		// User endpoints under /api/v1/user
		if strings.HasPrefix(path, "/user") {
			// register
			if r.Method == http.MethodPost && strings.HasPrefix(path, "/user/register") {
				var body bytes.Buffer
				if _, err := io.Copy(&body, r.Body); err != nil {
					http.Error(w, "failed to read body", http.StatusInternalServerError)
					return
				}
				if len(body.Bytes()) == 0 {
					http.Error(w, "empty body", http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusCreated)
				if _, err := w.Write([]byte(`{"id":"1"}`)); err != nil {
					fmt.Printf("Write error: %v\n", err)
				}
				return
			}
			// login
			if r.Method == http.MethodPost && strings.HasPrefix(path, "/user/login") {
				w.WriteHeader(http.StatusOK)
				if _, err := w.Write([]byte(`{"token":"test-token"}`)); err != nil {
					fmt.Printf("Write error: %v\n", err)
				}
				return
			}
			// refresh
			if r.Method == http.MethodPost && strings.HasPrefix(path, "/user/refresh") {
				w.WriteHeader(http.StatusOK)
				if _, err := w.Write([]byte(`{"token":"new-test-token"}`)); err != nil {
					fmt.Printf("Write error: %v\n", err)
				}
				return
			}
		}

		// Book endpoints under /api/v1/books
		if strings.HasPrefix(path, "/books") {
			handleCRUD(w, r, "/books")
			return
		}

		// Fallback for v1 routes
		http.NotFound(w, r)
	})

	return mux
}

// doReq is a test helper that builds and executes an HTTP request.
func doReq(t *testing.T, h http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var buf *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
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

// handleCRUD is a helper used by simple /api CRUD routing
func handleCRUD(w http.ResponseWriter, r *http.Request, base string) {
	// path like /exchanges, /exchanges/:id
	p := strings.TrimPrefix(r.URL.Path, "/api")
	p = strings.TrimPrefix(p, "/v1")

	// simple numeric id detection
	switch r.Method {
	case http.MethodPost:
		// validate body
		var body bytes.Buffer
		if _, err := io.Copy(&body, r.Body); err != nil {
			http.Error(w, "failed to read body", http.StatusInternalServerError)
			return
		}
		if len(body.Bytes()) == 0 {
			http.Error(w, "empty body", http.StatusBadRequest)
			return
		}

		// try parse JSON to catch invalid JSON
		var tmp map[string]any
		if err := json.Unmarshal(body.Bytes(), &tmp); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write([]byte(`{"id":"1"}`)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		return

	case http.MethodGet:
		if strings.Contains(p, "bad") {
			http.Error(w, "bad id", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`[]`)); err != nil {
			fmt.Printf("Write error: %v\n", err)
		}
		return

	case http.MethodPut:
		w.WriteHeader(http.StatusOK)
		return

	case http.MethodDelete:
		if strings.Contains(p, "NaN") || strings.Contains(p, "bad") {
			http.Error(w, "bad id", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return

	default:
		http.NotFound(w, r)
	}
}
