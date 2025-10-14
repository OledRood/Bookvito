package domain

// UserRepository defines methods for user data access
type UserRepository interface {
	Create(user *User) error
	GetByID(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
	List(limit, offset int) ([]*User, error)
}

// BookRepository defines methods for book data access
type BookRepository interface {
	Create(book *Book) error
	GetByID(id uint) (*Book, error)
	GetByOwnerID(ownerID uint) ([]*Book, error)
	Update(book *Book) error
	Delete(id uint) error
	List(limit, offset int) ([]*Book, error)
	Search(query string, limit, offset int) ([]*Book, error)
	GetAvailable(limit, offset int) ([]*Book, error)
}

// ExchangeRepository defines methods for exchange data access
type ExchangeRepository interface {
	Create(exchange *Exchange) error
	GetByID(id uint) (*Exchange, error)
	GetByRequesterID(requesterID uint) ([]*Exchange, error)
	GetByOwnerID(ownerID uint) ([]*Exchange, error)
	GetByBookID(bookID uint) ([]*Exchange, error)
	Update(exchange *Exchange) error
	Delete(id uint) error
	List(limit, offset int) ([]*Exchange, error)
}
