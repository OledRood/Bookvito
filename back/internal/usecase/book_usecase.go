package usecase

import (
	"bookvito/internal/domain"
	"errors"
	"time"

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

func (uc *BookUseCase) Request(bookID uuid.UUID, userID uuid.UUID) error {
	book, err := uc.bookRepo.GetByID(bookID)
	if err != nil {
		return err
	}
	if book.Status != domain.BookAvailable || book.Status == "" {
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

	expiresAt := time.Now().Add(48 * time.Hour) // Бронь истекает через 48 часов

	if err := uc.exchangeUseCaseRepo.Create(&domain.Exchange{
		UserID:    userID,
		BookID:    bookID,
		Status:    domain.ExchangeRequested,
		ExpiresAt: &expiresAt,
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
	exchanges, err := uc.exchangeUseCaseRepo.GetByBookID(bookID)
	if err != nil {
		return err
	}
	var requesterFound bool
	for _, ex := range exchanges {
		if ex.UserID == userID && ex.Status == domain.ExchangeRequested {
			requesterFound = true
			break
		}
	}
	if !requesterFound {
		return errors.New("only the user who requested the book can borrow it")
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

func (uc *BookUseCase) Return(updatedBook *domain.Book, userID uuid.UUID) error {
	if updatedBook.Title == "" {
		return errors.New("book title cannot be empty")
	}
	if updatedBook.Author == "" {
		return errors.New("book author cannot be empty")
	}
	println("Returning book status:", updatedBook.Status)

	bookFromDB, err := uc.bookRepo.GetByID(updatedBook.ID)
	if err != nil {
		return errors.New("Oooopsss... book not found in base")
	}

	if bookFromDB.Status != domain.BookBorrowed {
		return errors.New("only borrowed books can be returned")
	}

	// Обновляем только нужные поля у объекта, который мы получили из БД
	bookFromDB.Status = domain.BookAvailable
	bookFromDB.Title = updatedBook.Title
	bookFromDB.Author = updatedBook.Author
	bookFromDB.Description = updatedBook.Description
	bookFromDB.ImageURL = updatedBook.ImageURL
	bookFromDB.CurrentLocationID = updatedBook.CurrentLocationID
	bookFromDB.Condition = updatedBook.Condition // Обновляем состояние из запроса

	if err := uc.bookRepo.Update(bookFromDB); err != nil {
		return err
	}

	movement := &domain.BookMovementHistory{
		BookID:         bookFromDB.ID,
		UserID:         &userID,
		Action:         "returned",
		PreviousStatus: domain.BookBorrowed,
		NewStatus:      domain.BookAvailable,
		Notes:          "Book returned by user",
	}
	if err := uc.movementHistoryRepo.Create(movement); err != nil {
		return err
	}
	if err := uc.exchangeUseCaseRepo.Delete(bookFromDB.ID); err != nil {
		return err
	}

	return nil
}

func (uc *BookUseCase) DeleteBook(bookID, userID uuid.UUID) error {
	book, err := uc.bookRepo.GetByID(bookID)

	if err != nil {
		return errors.New("book not found")
	}

	// Новая проверка: книгу может удалить только тот, кто ее взял, и только если она в статусе "borrowed".
	if book.Status != domain.BookBorrowed {
		return errors.New("only a borrowed book can be marked as deleted by the borrower")
	}

	// Получаем последнюю запись в истории, чтобы проверить, кто взял книгу.
	history, err := uc.movementHistoryRepo.GetByBookID(bookID)
	if err != nil || len(history) == 0 {
		return errors.New("could not verify book history")
	}

	lastAction := history[0] // История отсортирована по убыванию даты
	if lastAction.Action != "borrowed" || lastAction.UserID == nil || *lastAction.UserID != userID {
		return errors.New("you are not the last user who borrowed this book")
	}

	// Старая проверка на владельца или админа (оставлена для справки)
	// if book.Owner.Role != domain.RoleAdmin && book.OwnerID != userID {
	// 	return errors.New("no rules to delete this book")
	// }

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
	if err := uc.movementHistoryRepo.Create(movement); err != nil {
		return err
	}
	if err := uc.exchangeUseCaseRepo.Delete(bookID); err != nil {
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

func (uc *BookUseCase) GetBookByID(bookID uuid.UUID) (*domain.Book, error) {
	return uc.bookRepo.GetByID(bookID)
}

func (uc *BookUseCase) GetBookMovementHistory(bookID uuid.UUID) ([]*domain.BookMovementHistory, error) {
	return uc.movementHistoryRepo.GetByBookID(bookID)
}
