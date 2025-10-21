package domain

import "github.com/google/uuid"

// UserUseCase интерфейс для работы с пользователями
type UserUseCase interface {
	RegisterUser(email, password, name string) (*TokenResponse, error)
	LoginUser(email, password string) (*TokenResponse, error)
	GetUserByID(id string) (*User, error)
	// UpdateUser(user *User) error
	// DeleteUser(id string) error
	// ListUsers(limit, offset int) ([]*User, error)
	RefreshToken(refreshToken string) (*TokenResponse, error)
	GetUserMovementHistory(userID string) ([]*BookMovementHistory, error)
}

// BookUseCase интерфейс для работы с книгами
type BookUseCase interface {
	CreateBook(book *Book) error
	GetSummaryBooksList() ([]*BookSummary, error)
	GetBooksList() ([]*Book, error)
	GetBookByID(id uuid.UUID) (*Book, error)
	DeleteBook(bookID uuid.UUID, userID uuid.UUID) error
	Request(bookID uuid.UUID, userID uuid.UUID) error
	Borrow(bookID uuid.UUID, userID uuid.UUID) error
	Return(bookID uuid.UUID, userID uuid.UUID) error

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
