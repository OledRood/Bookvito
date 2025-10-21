package domain

import "github.com/google/uuid"

// UserRepository defines methods for user data access
type UserRepository interface {
	Create(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]*User, error)
	GetByRefreshToken(refreshToken string) (*User, error)
}

// BookRepository defines methods for book data access
type BookRepository interface {
	Create(book *Book) error
	GetByID(id uuid.UUID) (*Book, error)
	Update(book *Book) error
	Delete(bookID uuid.UUID) error
	List(limit, offset int) ([]*Book, error)
	GetSummaryList(limit, offset int) ([]*BookSummary, error)
	Search(query string, limit, offset int) ([]*Book, error)
	GetByStatus(status BookStatus, limit, offset int) ([]*Book, error)
	GetByLocationID(locationID uuid.UUID) ([]*Book, error)
}

// ExchangeRepository defines methods for exchange data access
type ExchangeRepository interface {
	Create(exchange *Exchange) error
	GetByID(id uuid.UUID) (*Exchange, error)
	GetByUserID(userID uuid.UUID) ([]*Exchange, error)
	GetByBookID(bookID uuid.UUID) ([]*Exchange, error)
	Update(exchange *Exchange) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]*Exchange, error)
}

// LocationRepository defines methods for location data access
type LocationRepository interface {
	Create(location *Location) error
	GetByID(id uuid.UUID) (*Location, error)
	GetByAddress(address string) (*Location, error)
	GetAll() ([]Location, error)
	Update(location *Location) error
	Delete(id uuid.UUID) error
}

// ReviewRepository defines methods for review data access
type ReviewRepository interface {
	Create(review *Review) error
	GetByID(id uuid.UUID) (*Review, error)
	GetByBookID(bookID uuid.UUID) ([]Review, error)
	GetByUserID(userID uuid.UUID) ([]Review, error)
	Update(review *Review) error
	Delete(id uuid.UUID) error
}

// BookMovementHistoryRepository defines methods for book movement history data access
type BookMovementHistoryRepository interface {
	Create(movement *BookMovementHistory) error
	Update(movement *BookMovementHistory) error
	GetByID(id uuid.UUID) (*BookMovementHistory, error)
	GetByBookID(bookID uuid.UUID) ([]*BookMovementHistory, error)
	GetByExchangeID(exchangeID uuid.UUID) ([]*BookMovementHistory, error)
	GetByUserID(userID uuid.UUID) ([]*BookMovementHistory, error)
	List(limit, offset int) ([]*BookMovementHistory, error)
}
