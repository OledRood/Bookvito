package storage

import (
	"bookvito/internal/domain"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type LocalImageStorage struct {
	baseDir       string
	publicBaseURL string
}

func NewLocalImageStorage(baseDir, publicBaseURL string) (*LocalImageStorage, error) {
	normalizedBaseURL := normalizePublicBaseURL(publicBaseURL)
	if normalizedBaseURL == "" {
		return nil, fmt.Errorf("public base URL is required")
	}
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return nil, err
	}

	return &LocalImageStorage{
		baseDir:       baseDir,
		publicBaseURL: normalizedBaseURL,
	}, nil
}

func (s *LocalImageStorage) Save(filename, _ string, reader io.Reader) (string, error) {
	safeFilename, err := sanitizeFilename(filename)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(s.baseDir, safeFilename)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return "", err
	}

	return s.PublicURL(safeFilename), nil
}

func (s *LocalImageStorage) Exists(publicURL string) (bool, error) {
	filePath, err := s.pathFromURL(publicURL)
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *LocalImageStorage) Delete(publicURL string) error {
	filePath, err := s.pathFromURL(publicURL)
	if err != nil {
		return err
	}

	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *LocalImageStorage) PublicURL(filename string) string {
	return s.publicBaseURL + strings.TrimLeft(filename, "/")
}

func (s *LocalImageStorage) pathFromURL(publicURL string) (string, error) {
	value := strings.TrimSpace(publicURL)
	if value == "" {
		return "", domain.NewValidationError("image_url пустой")
	}
	if !strings.HasPrefix(value, s.publicBaseURL) {
		return "", domain.NewValidationError("image_url не принадлежит image storage")
	}

	filename := strings.TrimPrefix(value, s.publicBaseURL)
	filename, err := sanitizeFilename(filename)
	if err != nil {
		return "", err
	}

	return filepath.Join(s.baseDir, filename), nil
}

func normalizePublicBaseURL(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	if !strings.HasSuffix(trimmed, "/") {
		trimmed += "/"
	}
	return trimmed
}

func sanitizeFilename(filename string) (string, error) {
	value := strings.TrimSpace(filename)
	value = filepath.Base(value)
	if value == "" || value == "." {
		return "", domain.NewValidationError("Некорректный filename")
	}
	return value, nil
}
