package httptestpkg

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type bookCreate struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

type bookSearch struct {
	Query string `json:"query"`
}

func TestBook_Fullflow_Create_Get_Search_Update_Delete(t *testing.T) {
	r := NewRouter()

	// Create valid
	b, err := json.Marshal(bookCreate{Title: "Dune", Author: "Frank Herbert"})
	if err != nil {
		t.Fatalf("failed to marshal create request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/books", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	if rw.Code != http.StatusCreated && rw.Code != http.StatusOK {
		t.Fatalf("create expected 201/200, got %d, body=%s", rw.Code, rw.Body.String())
	}

	var created map[string]any
	if err := json.Unmarshal(rw.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to unmarshal create response: %v", err)
	}

	id, _ := created["id"].(string)
	if id == "" {
		// try numeric id too
		if _, ok := created["id"].(float64); !ok {
			t.Fatalf("expected id in create response")
		}
	}

	// Get by ID
	req = httptest.NewRequest(http.MethodGet, "/api/v1/books/1", nil)
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("get by id expected 200, got %d", rw.Code)
	}

	// Search
	b, err = json.Marshal(bookSearch{Query: "Dune"})
	if err != nil {
		t.Fatalf("failed to marshal search request: %v", err)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/v1/books/search", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("search expected 200, got %d", rw.Code)
	}

	// Update
	upd := map[string]any{"title": "Dune (Updated)"}
	b, err = json.Marshal(upd)
	if err != nil {
		t.Fatalf("failed to marshal update request: %v", err)
	}

	req = httptest.NewRequest(http.MethodPut, "/api/v1/books/1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("update expected 200, got %d", rw.Code)
	}

	// Delete
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/books/1", nil)
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK && rw.Code != http.StatusNoContent {
		t.Fatalf("delete expected 200/204, got %d", rw.Code)
	}
}

func TestBook_List_Filter_Pagination_InvalidInputs(t *testing.T) {
	r := NewRouter()

	// List default
	req := httptest.NewRequest(http.MethodGet, "/api/v1/books", nil)
	rw := httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("list expected 200, got %d", rw.Code)
	}

	// List with filters
	req = httptest.NewRequest(http.MethodGet, "/api/v1/books?author=Herbert&title=Dune&page=1&size=10", nil)
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("filtered list expected 200, got %d", rw.Code)
	}

	// Invalid page/size
	req = httptest.NewRequest(http.MethodGet, "/api/v1/books?page=-1&size=0", nil)
	rw = httptest.NewRecorder()
	r.ServeHTTP(rw, req)

	if rw.Code == http.StatusOK {
		t.Fatalf("invalid pagination should not be 200")
	}
}
