package domain

// BookMeta represents essential metadata about a book returned from external providers.
type BookMeta struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

// BookMetadataProvider describes services capable of fetching book metadata by query.
type BookMetadataProvider interface {
	Search(query string) (*BookMeta, error)
}
