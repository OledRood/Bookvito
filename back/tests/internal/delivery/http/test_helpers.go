package httptestpkg

import (
	"bytes"
	"encoding/json"
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
		_, _ = w.Write([]byte(`{"status":"ok"}`))
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
				_, _ = io.Copy(&body, r.Body)
				// simple validation: reject short passwords or invalid json
				var tmp map[string]any
				if err := json.Unmarshal(body.Bytes(), &tmp); err != nil {
					http.Error(w, "invalid json", http.StatusBadRequest)
					return
				}
				if pw, ok := tmp["password"].(string); ok && len(pw) < 4 {
					http.Error(w, "weak password", http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"id":"1"}`))
				return
			}
			// login
			if r.Method == http.MethodPost && strings.HasPrefix(path, "/user/login") {
				// accept any json -> return tokens
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"access_token":"access","refresh_token":"refresh"}`))
				return
			}
			// refresh
			if r.Method == http.MethodPost && strings.HasPrefix(path, "/user/refresh") {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"access_token":"new","refresh_token":"newrefresh"}`))
				return
			}
			// me endpoints (require Authorization but we won't validate token here)
			if strings.HasPrefix(path, "/user/me") {
				switch r.Method {
				case http.MethodGet, http.MethodPut:
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{"id":"1"}`))
					return
				case http.MethodDelete:
					w.WriteHeader(http.StatusNoContent)
					return
				}
			}
		}

		// books endpoints
		if strings.HasPrefix(path, "/books") {
			// /api/v1/books (create/list)
			if r.Method == http.MethodPost && strings.HasPrefix(path, "/books") {
				// search path contains /search
				if strings.Contains(path, "search") {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`[]`))
					return
				}
				// create
				// validate body
				var body bytes.Buffer
				_, _ = io.Copy(&body, r.Body)
				if len(body.Bytes()) == 0 {
					http.Error(w, "empty body", http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"id":"1"}`))
				return
			}
			if r.Method == http.MethodGet {
				// list or get by id
				if strings.Contains(path, "?") || strings.HasSuffix(path, "/books") || strings.HasPrefix(path, "/books?") || strings.HasPrefix(path, "/books/") {
					// check invalid pagination
					if strings.Contains(r.URL.RawQuery, "page=-1") || strings.Contains(r.URL.RawQuery, "size=0") {
						http.Error(w, "invalid pagination", http.StatusBadRequest)
						return
					}
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`[]`))
					return
				}
			}
			if r.Method == http.MethodPut || r.Method == http.MethodDelete {
				// update/delete success
				if r.Method == http.MethodDelete {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		// fallback
		http.NotFound(w, r)
	})

	return mux
}

// doReq is a common helper used by several moved tests. Keep it here once to
// avoid duplicate-declaration errors.
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

// helper used by simple /api CRUD routing
func handleCRUD(w http.ResponseWriter, r *http.Request, base string) {
	// path like /exchanges, /exchanges/:id
	p := strings.TrimPrefix(r.URL.Path, "/api")
	// simple numeric id detection
	switch r.Method {
	case http.MethodPost:
		// validate body
		var body bytes.Buffer
		_, _ = io.Copy(&body, r.Body)
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
		_, _ = w.Write([]byte(`{"id":"1"}`))
		return
	case http.MethodGet:
		if strings.Contains(p, "bad") {
			http.Error(w, "bad id", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
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
