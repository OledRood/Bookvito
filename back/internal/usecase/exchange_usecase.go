package usecase

import (
	"bookvito/internal/domain"
	"errors"
)

type ExchangeUseCase struct {
	exchangeRepo domain.ExchangeRepository
	bookRepo     domain.BookRepository
	userRepo     domain.UserRepository
}

// NewExchangeUseCase creates a new exchange use case
func NewExchangeUseCase(exchangeRepo domain.ExchangeRepository, bookRepo domain.BookRepository, userRepo domain.UserRepository) *ExchangeUseCase {
	return &ExchangeUseCase{
		exchangeRepo: exchangeRepo,
		bookRepo:     bookRepo,
		userRepo:     userRepo,
	}
}

// CreateExchangeRequest creates a new exchange request
func (uc *ExchangeUseCase) CreateExchangeRequest(requesterID, bookID uint, message string) (*domain.Exchange, error) {
	// Get the book
	book, err := uc.bookRepo.GetByID(bookID)
	if err != nil {
		return nil, err
	}

	// Check if the book is available
	if !book.Available {
		return nil, errors.New("book is not available for exchange")
	}

	// Check if requester is not the owner
	if book.OwnerID == requesterID {
		return nil, errors.New("you cannot request your own book")
	}

	exchange := &domain.Exchange{
		RequesterID: requesterID,
		OwnerID:     book.OwnerID,
		BookID:      bookID,
		Status:      "pending",
		Message:     message,
	}

	if err := uc.exchangeRepo.Create(exchange); err != nil {
		return nil, err
	}

	return exchange, nil
}

// GetExchangeByID retrieves an exchange by ID
func (uc *ExchangeUseCase) GetExchangeByID(id uint) (*domain.Exchange, error) {
	return uc.exchangeRepo.GetByID(id)
}

// GetExchangesByRequester retrieves all exchanges requested by a user
func (uc *ExchangeUseCase) GetExchangesByRequester(requesterID uint) ([]*domain.Exchange, error) {
	return uc.exchangeRepo.GetByRequesterID(requesterID)
}

// GetExchangesByOwner retrieves all exchanges for books owned by a user
func (uc *ExchangeUseCase) GetExchangesByOwner(ownerID uint) ([]*domain.Exchange, error) {
	return uc.exchangeRepo.GetByOwnerID(ownerID)
}

// AcceptExchange accepts an exchange request
func (uc *ExchangeUseCase) AcceptExchange(exchangeID, ownerID uint) error {
	exchange, err := uc.exchangeRepo.GetByID(exchangeID)
	if err != nil {
		return err
	}

	// Check if the user is the owner
	if exchange.OwnerID != ownerID {
		return errors.New("unauthorized: only the owner can accept this exchange")
	}

	// Check if exchange is pending
	if exchange.Status != "pending" {
		return errors.New("exchange is not in pending status")
	}

	exchange.Status = "accepted"
	if err := uc.exchangeRepo.Update(exchange); err != nil {
		return err
	}

	// Mark the book as unavailable
	book, err := uc.bookRepo.GetByID(exchange.BookID)
	if err != nil {
		return err
	}

	book.Available = false
	return uc.bookRepo.Update(book)
}

// RejectExchange rejects an exchange request
func (uc *ExchangeUseCase) RejectExchange(exchangeID, ownerID uint) error {
	exchange, err := uc.exchangeRepo.GetByID(exchangeID)
	if err != nil {
		return err
	}

	// Check if the user is the owner
	if exchange.OwnerID != ownerID {
		return errors.New("unauthorized: only the owner can reject this exchange")
	}

	// Check if exchange is pending
	if exchange.Status != "pending" {
		return errors.New("exchange is not in pending status")
	}

	exchange.Status = "rejected"
	return uc.exchangeRepo.Update(exchange)
}

// CompleteExchange marks an exchange as completed
func (uc *ExchangeUseCase) CompleteExchange(exchangeID, ownerID uint) error {
	exchange, err := uc.exchangeRepo.GetByID(exchangeID)
	if err != nil {
		return err
	}

	// Check if the user is the owner
	if exchange.OwnerID != ownerID {
		return errors.New("unauthorized: only the owner can complete this exchange")
	}

	// Check if exchange is accepted
	if exchange.Status != "accepted" {
		return errors.New("exchange must be accepted before it can be completed")
	}

	exchange.Status = "completed"
	return uc.exchangeRepo.Update(exchange)
}

// ListExchanges retrieves a list of exchanges
func (uc *ExchangeUseCase) ListExchanges(limit, offset int) ([]*domain.Exchange, error) {
	return uc.exchangeRepo.List(limit, offset)
}
