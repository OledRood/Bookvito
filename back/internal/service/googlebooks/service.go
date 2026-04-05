package googlebooks

import (
	"bookvito/internal/domain"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GoogleBooksResponse mirrors the minimal subset of the Google Books API response we need.
type GoogleBooksResponse struct {
	Items []struct {
		VolumeInfo struct {
			Title       string   `json:"title"`
			Authors     []string `json:"authors"`
			Description string   `json:"description"`
			ImageLinks  struct {
				Thumbnail string `json:"thumbnail"`
			} `json:"imageLinks"`
		} `json:"volumeInfo"`
	} `json:"items"`
}

// Service implements domain.BookMetadataProvider using Google Books API.
type Service struct {
	apiKey string
	client *http.Client
}

// New creates a Google Books service with the provided API key.
func New(apiKey string) *Service {
	return &Service{
		apiKey: apiKey,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

// Search finds book metadata by query using Google Books.
func (s *Service) Search(query string) (*domain.BookMeta, error) {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return nil, domain.NewValidationError("q is required")
	}
	if s.apiKey == "" {
		return nil, errors.New("google books api key is not configured")
	}

	encodedQuery := url.QueryEscape(trimmed)
	requestURL := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=%s&key=%s", encodedQuery, s.apiKey)

	resp, err := s.client.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google books responded with status %d", resp.StatusCode)
	}

	var result GoogleBooksResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, domain.NewNotFoundError("book not found")
	}

	info := result.Items[0].VolumeInfo

	return &domain.BookMeta{
		Title:       strings.TrimSpace(info.Title),
		Author:      strings.TrimSpace(strings.Join(info.Authors, ", ")),
		Description: strings.TrimSpace(info.Description),
		ImageURL:    strings.TrimSpace(info.ImageLinks.Thumbnail),
	}, nil
}

// Ensure Service implements the provider interface at compile time.
var _ domain.BookMetadataProvider = (*Service)(nil)
