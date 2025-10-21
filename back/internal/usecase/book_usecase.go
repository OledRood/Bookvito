package usecase

import (
	"bookvito/internal/domain"
	"errors"

	"github.com/google/uuid"
)

type BookUseCase struct {
	bookRepo            domain.BookRepository
	movementHistoryRepo domain.BookMovementHistoryRepository
	exchangeUseCaseRepo domain.ExchangeRepository
}

func NewBookUseCase(bookRepo domain.BookRepository, movementHistoryRepo domain.BookMovementHistoryRepository, exchangeUseCaseRepo domain.ExchangeRepository) *BookUseCase {
	return &BookUseCase{
		bookRepo:            bookRepo,
		movementHistoryRepo: movementHistoryRepo,
		exchangeUseCaseRepo: exchangeUseCaseRepo,
	}
}

// CreateBook creates a new book
func (uc *BookUseCase) CreateBook(book *domain.Book) error {

	if err := uc.bookRepo.Create(book); err != nil {
		return err
	}

	// Создаем запись в истории перемещений
	movement := &domain.BookMovementHistory{
		BookID:         book.ID,
		ToLocationID:   book.CurrentLocationID,
		UserID:         &book.OwnerID,
		Action:         "created",
		Notes:          "Книга добавлена в систему",
		PreviousStatus: "",
		NewStatus:      domain.BookAvailable,
	}

	if err := uc.movementHistoryRepo.Create(movement); err != nil {
		return err
	}

	return nil
}

func (uc *BookUseCase) GetSummaryBooksList() ([]*domain.BookSummary, error) {
	return uc.bookRepo.GetSummaryList(100, 0)
}

func (uc *BookUseCase) GetBooksList() ([]*domain.Book, error) {
	return uc.bookRepo.List(100, 0)
}

func (uc *BookUseCase) GetBookByID(id uuid.UUID) (*domain.Book, error) {
	return uc.bookRepo.GetByID(id)
}

func (uc *BookUseCase) DeleteBook(bookID, userID uuid.UUID) error {
	book, err := uc.bookRepo.GetByID(bookID)
	if err != nil {
		return err
	}

	previousStatus := book.Status

	book.Status = domain.BookDeleted

	if err := uc.bookRepo.Update(book); err != nil {
		return err
	}

	// 5. Создаем запись в истории об этом событии
	movement := &domain.BookMovementHistory{
		BookID:         bookID,
		UserID:         &userID, // Пользователь, который выполнил действие
		Action:         "deleted",
		PreviousStatus: previousStatus,
		NewStatus:      domain.BookDeleted,
	}

	// 6. Сохраняем запись в истории
	return uc.movementHistoryRepo.Create(movement)
}

func (uc *BookUseCase) GetBookMovementHistory(bookID uuid.UUID) ([]*domain.BookMovementHistory, error) {
	return uc.movementHistoryRepo.GetByBookID(bookID)
}

func (uc *BookUseCase) Request(bookID uuid.UUID, userID uuid.UUID) error {
	book, err := uc.bookRepo.GetByID(bookID)
	if err != nil {
		return err
	}
	if book.Status != domain.BookAvailable {
		return errors.New("book is not available for request")
	}

	book.Status = domain.BookRequested
	if err := uc.bookRepo.Update(book); err != nil {
		return err
	}

	// Создаем запись в истории перемещений
	movement := &domain.BookMovementHistory{
		BookID:         book.ID,
		UserID:         &userID,
		Action:         "requested",
		PreviousStatus: domain.BookAvailable,
		NewStatus:      domain.BookRequested,
		Notes:          "Book requested by user",
	}
	if err := uc.movementHistoryRepo.Create(movement); err != nil {
		return err
	}

	if err := uc.exchangeUseCaseRepo.Create(&domain.Exchange{
		UserID: userID,
		BookID: bookID,
		Status: domain.ExchangeRequested,
	}); err != nil {
		return err
	}

	return nil
}

func (uc *BookUseCase) Borrow(bookID uuid.UUID, userID uuid.UUID) error {
	book, err := uc.bookRepo.GetByID(bookID)
	if err != nil {
		return err
	}
	if book.Status != domain.BookRequested {
		return errors.New("book is not available for borrowing")
	}

	book.Status = domain.BookBorrowed
	if err := uc.bookRepo.Update(book); err != nil {
		return err
	}

	// Создаем запись в истории перемещений
	movement := &domain.BookMovementHistory{
		BookID:         book.ID,
		UserID:         &userID,
		Action:         "borrowed",
		PreviousStatus: domain.BookAvailable,
		NewStatus:      domain.BookBorrowed,
		Notes:          "Book borrowed by user",
	}
	if err := uc.movementHistoryRepo.Create(movement); err != nil {
		return err
	}

	return nil
}

func (uc *BookUseCase) Return(bookID uuid.UUID, userID uuid.UUID) error {
	book, err := uc.bookRepo.GetByID(bookID)
	if err != nil {
		return err
	}
	if book.Status != domain.BookBorrowed {
		return errors.New("book is not currently borrowed")
	}

	book.Status = domain.BookAvailable
	if err := uc.bookRepo.Update(book); err != nil {
		return err
	}

	// Создаем запись в истории перемещений
	movement := &domain.BookMovementHistory{
		BookID:         book.ID,
		UserID:         &userID,
		Action:         "returned",
		PreviousStatus: domain.BookBorrowed,
		NewStatus:      domain.BookAvailable,
		Notes:          "Book returned by user",
	}
	if err := uc.movementHistoryRepo.Create(movement); err != nil {
		return err
	}
	if err := uc.exchangeUseCaseRepo.Delete(bookID); err != nil {
		return err
	}

	return nil
}
