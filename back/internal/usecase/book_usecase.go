package usecase

import (
	"bookvito/internal/domain"
	"errors"
)

type BookUseCase struct {
	bookRepo domain.BookRepository
}

// NewBookUseCase creates a new book use case
func NewBookUseCase(bookRepo domain.BookRepository) *BookUseCase {
	return &BookUseCase{
		bookRepo: bookRepo,
	}
}

// CreateBook creates a new book
func (uc *BookUseCase) CreateBook(title, author, isbn, description, condition string, ownerID uint) (*domain.Book, error) {
	book := &domain.Book{
		Title:       title,
		Author:      author,
		ISBN:        isbn,
		Description: description,
		Condition:   condition,
		OwnerID:     ownerID,
		Available:   true,
	}

	if err := uc.bookRepo.Create(book); err != nil {
		return nil, err
	}

	return book, nil
}

// GetBookByID retrieves a book by ID
func (uc *BookUseCase) GetBookByID(id uint) (*domain.Book, error) {
	return uc.bookRepo.GetByID(id)
}

// GetBooksByOwner retrieves all books owned by a user
func (uc *BookUseCase) GetBooksByOwner(ownerID uint) ([]*domain.Book, error) {
	return uc.bookRepo.GetByOwnerID(ownerID)
}

// UpdateBook updates book information
func (uc *BookUseCase) UpdateBook(book *domain.Book) error {
	return uc.bookRepo.Update(book)
}

// DeleteBook deletes a book
func (uc *BookUseCase) DeleteBook(id uint, ownerID uint) error {
	book, err := uc.bookRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Check if the user is the owner
	if book.OwnerID != ownerID {
		return errors.New("unauthorized: you can only delete your own books")
	}

	return uc.bookRepo.Delete(id)
}

// ListBooks retrieves a list of books
func (uc *BookUseCase) ListBooks(limit, offset int) ([]*domain.Book, error) {
	return uc.bookRepo.List(limit, offset)
}

// SearchBooks searches for books by title, author, or description
func (uc *BookUseCase) SearchBooks(query string, limit, offset int) ([]*domain.Book, error) {
	return uc.bookRepo.Search(query, limit, offset)
}

// GetAvailableBooks retrieves all available books
func (uc *BookUseCase) GetAvailableBooks(limit, offset int) ([]*domain.Book, error) {
	return uc.bookRepo.GetAvailable(limit, offset)
}

// MarkBookAsUnavailable marks a book as unavailable
func (uc *BookUseCase) MarkBookAsUnavailable(id uint) error {
	book, err := uc.bookRepo.GetByID(id)
	if err != nil {
		return err
	}

	book.Available = false
	return uc.bookRepo.Update(book)
}

// MarkBookAsAvailable marks a book as available
func (uc *BookUseCase) MarkBookAsAvailable(id uint) error {
	book, err := uc.bookRepo.GetByID(id)
	if err != nil {
		return err
	}

	book.Available = true
	return uc.bookRepo.Update(book)
}
