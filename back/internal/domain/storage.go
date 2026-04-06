package domain

import "io"

// ImageStorage provides a storage-agnostic API for book images.
type ImageStorage interface {
	Save(filename, contentType string, reader io.Reader) (string, error)
	Exists(publicURL string) (bool, error)
	Delete(publicURL string) error
	PublicURL(filename string) string
}
