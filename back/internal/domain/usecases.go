package domain

import "github.com/google/uuid"

// UserUseCase интерфейс для работы с пользователями
type UserUseCase interface {
	RegisterUser(email, password, name string) (*TokenResponse, error)
	LoginUser(email, password string) (*TokenResponse, error)
	GetUserByID(id string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id string) error
	// ListUsers(limit, offset int) ([]*User, error)
	RefreshToken(refreshToken string) (*TokenResponse, error)
	GetUserMovementHistory(userID string) ([]*BookMovementHistory, error)
}

// BookUseCase интерфейс для работы с книгами
type BookUseCase interface {
	CreateBook(book *Book) error
	// GetSummaryBooksList returns public summaries. If userID is non-nil,
	// summaries for books owned by or interacted with by that user will be excluded.
	GetSummaryBooksList(userID *uuid.UUID) ([]*BookSummary, error)
	GetBooksList() ([]*Book, error)
	GetBookByID(bookID uuid.UUID) (*Book, error)
	// SearchBooks performs a text search over books. If userID is non-nil,
	// search results should follow the same visibility rules as GetSummaryBooksList
	// (exclude user's own books and books the user has returned).
	SearchBooks(query string, limit, offset int, userID *uuid.UUID) ([]*BookSummary, error)
	// DeleteBook deletes (marks deleted) a book. If isAdmin is true the caller is allowed
	// to delete any book. Otherwise the caller must be the user who currently has the
	// book on their shelf (i.e. borrowed by them).
	DeleteBook(bookID uuid.UUID, userID uuid.UUID, isAdmin bool, reason string) error
	// Request creates a reservation (exchange) for a book. Optional locationID
	// is the chosen pickup location; if nil the usecase may fall back to the
	// book's current location.
	Request(bookID uuid.UUID, userID uuid.UUID, locationID *uuid.UUID) error
	Borrow(bookID uuid.UUID, userID uuid.UUID) error
	Return(updatedBook *Book, userID uuid.UUID) error
	// SetBookImage allows setting or updating the public image URL for a book
	SetBookImage(bookID uuid.UUID, imageURL string) error

	// GetBookByID(id uuid.UUID) (*Book, error)
	// UpdateBook(book *Book) error
	// DeleteBook(id uuid.UUID) error
	// ListBooks(limit, offset int) ([]*Book, error)
	// SearchBooks(query string) ([]*Book, error)
	// GetBooksByOwner(ownerID uuid.UUID) ([]*Book, error)
	// GetBooksByLocation(locationID uuid.UUID) ([]*Book, error)
	// GetAvailableBooks() ([]*Book, error)

	// Методы для работы с историей перемещений
	GetBookMovementHistory(bookID uuid.UUID) ([]*BookMovementHistory, error)
	// Additional convenience methods used by handlers
	GetBooksByOwner(ownerID uuid.UUID) ([]*Book, error)
	// Get books related to exchanges for a user filtered by exchange status (requested/borrowed/returned)
	GetBooksFromExchangesByStatus(userID uuid.UUID, status ExchangeStatus) ([]*Book, error)
	// Get exchanges (with book and location preloaded) for a user filtered by status
	GetExchangesFromUserByStatus(userID uuid.UUID, status ExchangeStatus) ([]*Exchange, error)

	// Get aggregated statistics for authenticated user's books
	GetMyBooksStats(ownerID uuid.UUID) (*MyBooksStats, error)

	// Reservation management: extend and cancel
	ExtendReservation(bookID uuid.UUID, userID uuid.UUID) error
	CancelReservation(bookID uuid.UUID, userID uuid.UUID) error

	// Get statistics for a single book (by book ID)
	GetBookStats(bookID uuid.UUID) (*MyBooksStats, error)
}

// MyBooksStats represents aggregated statistics for a user's books
type MyBooksStats struct {
	TotalBooks           int64            `json:"total_books"`
	StatusCounts         map[string]int64 `json:"status_counts"`
	ConditionCounts      map[string]int64 `json:"condition_counts"`
	TotalUniqueBorrowers int64            `json:"total_unique_borrowers"`
	FirstLocation        *string          `json:"first_location,omitempty"`
	LastLocation         *string          `json:"last_location,omitempty"`
	// Current status/condition (useful for per-book stats)
	CurrentStatus    *string `json:"current_status,omitempty"`
	CurrentCondition *string `json:"current_condition,omitempty"`
}

// ExchangeUseCase интерфейс для работы с обменом книг
type ExchangeUseCase interface {
	// CreateExchangeRequest(userID, bookID uuid.UUID, locationID *uuid.UUID) (*Exchange, error)
	// GetExchangeByID(id uuid.UUID) (*Exchange, error)
	// GetExchangesByUser(userID uuid.UUID) ([]*Exchange, error)
	// GetExchangesByBook(bookID uuid.UUID) ([]*Exchange, error)
	// ApproveExchange(exchangeID uuid.UUID, dueDate *time.Time) error
	// BorrowBook(exchangeID uuid.UUID) error
	// ReturnBook(exchangeID uuid.UUID) error
	// CancelExchange(exchangeID uuid.UUID) error
	// ListExchanges(limit, offset int) ([]*Exchange, error)
}

type LocationUseCase interface {
	Create(location *Location) error
	GetByID(id uuid.UUID) (*Location, error)
	GetAll() ([]Location, error)
	Update(location *Location) error
	Delete(id uuid.UUID) error
}

// TokenResponse структура ответа с токенами
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
