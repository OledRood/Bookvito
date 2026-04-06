//go:build integration

package integration_test

import (
	"bookvito/internal/domain"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserEndpoints_RegisterDuplicateConflictAndRefresh(t *testing.T) {
	cfg := testConfig()
	_, _, router := setupIntegrationTest(t, cfg)

	registerBody := map[string]string{
		"email":    "reader@example.com",
		"password": "password123",
		"name":     "Reader",
	}

	first := performJSON(t, router, http.MethodPost, "/api/v1/users/registration", registerBody, "")
	assertStatus(t, first, http.StatusCreated)
	var firstPayload struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	firstPayload = decodeJSON[struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}](t, first.Body)
	if firstPayload.AccessToken == "" || firstPayload.RefreshToken == "" {
		t.Fatalf("expected tokens in register response: %+v", firstPayload)
	}

	duplicate := performJSON(t, router, http.MethodPost, "/api/v1/users/registration", registerBody, "")
	assertStatus(t, duplicate, http.StatusConflict)
	duplicatePayload := decodeJSON[map[string]any](t, duplicate.Body)
	if duplicatePayload["code"] != string(domain.ErrorCodeConflict) {
		t.Fatalf("expected conflict code, got %+v", duplicatePayload)
	}

	refresh := performJSON(t, router, http.MethodPost, "/api/v1/users/refresh", map[string]string{
		"refresh_token": firstPayload.RefreshToken,
	}, "")
	assertStatus(t, refresh, http.StatusOK)
}

func TestBookEndpoints_UploadCreateAndValidateListParams(t *testing.T) {
	cfg := testConfig()
	db, _, router := setupIntegrationTest(t, cfg)

	createUserWithRole(t, db, "owner@example.com", "password123", string(domain.RoleUser))
	token := authToken(t, router, "owner@example.com", "password123")

	upload := performMultipartImageUpload(t, router, "/api/v1/books/upload", token)
	assertStatus(t, upload, http.StatusOK)
	uploadPayload := decodeJSON[map[string]string](t, upload.Body)
	imageURL := uploadPayload["url"]
	if imageURL == "" {
		t.Fatalf("expected uploaded image URL")
	}

	create := performJSON(t, router, http.MethodPost, "/api/v1/books/create", map[string]any{
		"title":               "Distributed Systems",
		"author":              "Tanenbaum",
		"description":         "Classic",
		"condition":           "good",
		"image_url":           imageURL,
		"current_location_id": nil,
	}, token)
	assertStatus(t, create, http.StatusCreated)

	list := performJSON(t, router, http.MethodGet, "/api/v1/books/list?limit=abc", nil, "")
	assertStatus(t, list, http.StatusBadRequest)
	listPayload := decodeJSON[map[string]any](t, list.Body)
	if listPayload["code"] != string(domain.ErrorCodeValidation) {
		t.Fatalf("expected validation code, got %+v", listPayload)
	}
}

func TestAdminEndpoint_RequiresRoleAndAllowsAdminUpdate(t *testing.T) {
	cfg := testConfig()
	db, _, router := setupIntegrationTest(t, cfg)

	createUserWithRole(t, db, "admin@example.com", "password123", string(domain.RoleAdmin))
	createUserWithRole(t, db, "user@example.com", "password123", string(domain.RoleUser))
	target := createUserWithRole(t, db, "target@example.com", "password123", string(domain.RoleUser))

	userToken := authToken(t, router, "user@example.com", "password123")
	adminToken := authToken(t, router, "admin@example.com", "password123")

	forbidden := performJSON(t, router, http.MethodGet, "/api/v1/admin/users", nil, userToken)
	assertStatus(t, forbidden, http.StatusForbidden)
	forbiddenPayload := decodeJSON[map[string]any](t, forbidden.Body)
	if forbiddenPayload["code"] != string(domain.ErrorCodeForbidden) {
		t.Fatalf("expected forbidden code, got %+v", forbiddenPayload)
	}

	update := performJSON(t, router, http.MethodPut, fmt.Sprintf("/api/v1/admin/users/%s/role", target.ID), map[string]string{
		"role": "moder",
	}, adminToken)
	assertStatus(t, update, http.StatusOK)
}

func TestBookAutoFill_UsesConfiguredBaseURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"items":[{"volumeInfo":{"title":"Stubbed Book","authors":["Stub Author"],"description":"Stub desc","imageLinks":{"thumbnail":"http://example.com/thumb.png"}}}]}`)
	}))
	defer server.Close()

	cfg := testConfig()
	cfg.GoogleBooksBaseURL = server.URL
	_, _, router := setupIntegrationTest(t, cfg)

	response := performJSON(t, router, http.MethodGet, "/api/v1/books/auto-fill?q=isbn:123", nil, "")
	assertStatus(t, response, http.StatusOK)
	payload := decodeJSON[map[string]any](t, response.Body)
	if payload["title"] != "Stubbed Book" {
		t.Fatalf("unexpected auto-fill payload: %+v", payload)
	}
}

func TestBookAutoFill_HandlesProviderFailures(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		w.Header().Set("Content-Type", "application/json")
		switch query {
		case "stub-empty":
			fmt.Fprint(w, `{"items":[]}`)
		case "stub-400":
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"error":"bad request"}`)
		case "stub-500":
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"error":"internal"}`)
		default:
			fmt.Fprint(w, `{"items":[{"volumeInfo":{"title":"Stubbed Book","authors":["Stub Author"],"description":"Stub desc","imageLinks":{"thumbnail":"http://example.com/thumb.png"}}}]}`)
		}
	}))
	defer server.Close()

	cfg := testConfig()
	cfg.GoogleBooksBaseURL = server.URL
	_, _, router := setupIntegrationTest(t, cfg)

	emptyResponse := performJSON(t, router, http.MethodGet, "/api/v1/books/auto-fill?q=stub-empty", nil, "")
	assertStatus(t, emptyResponse, http.StatusNotFound)
	emptyPayload := decodeJSON[map[string]any](t, emptyResponse.Body)
	if emptyPayload["code"] != string(domain.ErrorCodeNotFound) {
		t.Fatalf("expected not_found code, got %+v", emptyPayload)
	}

	for _, query := range []string{"stub-400", "stub-500"} {
		response := performJSON(t, router, http.MethodGet, fmt.Sprintf("/api/v1/books/auto-fill?q=%s", query), nil, "")
		assertStatus(t, response, http.StatusInternalServerError)
		payload := decodeJSON[map[string]any](t, response.Body)
		if payload["code"] != string(domain.ErrorCodeInternal) {
			t.Fatalf("expected internal code for %s, got %+v", query, payload)
		}
	}
}
