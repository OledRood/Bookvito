package usecase

import (
	"bookvito/internal/domain"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookUseCase struct {
	bookRepo            domain.BookRepository
	movementHistoryRepo domain.BookMovementHistoryRepository
	exchangeUseCaseRepo domain.ExchangeRepository
	locationRepo        domain.LocationRepository
}

func NewBookUseCase(bookRepo domain.BookRepository, movementHistoryRepo domain.BookMovementHistoryRepository, exchangeUseCaseRepo domain.ExchangeRepository, locationRepo domain.LocationRepository) *BookUseCase {
	return &BookUseCase{
		bookRepo:            bookRepo,
		movementHistoryRepo: movementHistoryRepo,
		exchangeUseCaseRepo: exchangeUseCaseRepo,
		locationRepo:        locationRepo,
	}
}

func (uc *BookUseCase) CreateBook(book *domain.Book) error {
	if err := uc.validateBookForCreate(book); err != nil {
		return err
	}

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

func (uc *BookUseCase) Request(bookID uuid.UUID, userID uuid.UUID, locationID *uuid.UUID) error {
	if err := validateBookID(bookID); err != nil {
		return err
	}
	if userID == uuid.Nil {
		return domain.NewValidationError("Неверный идентификатор пользователя")
	}
	if err := uc.validateLocation(locationID); err != nil {
		return err
	}

	book, err := uc.bookRepo.GetByID(bookID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.NewNotFoundError("Книга не найдена")
		}
		return err
	}
	if book.OwnerID == userID {
		return domain.NewForbiddenError("Нельзя забронировать собственную книгу")
	}
	if book.Status != domain.BookAvailable || book.Status == "" {
		return domain.NewConflictError(fmt.Sprintf("Книга недоступна для бронирования (статус=%s)", book.Status))
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

	// prefer explicit locationID from the request; otherwise fall back to book's current location
	chosenLocationID := locationID
	if chosenLocationID == nil {
		chosenLocationID = book.CurrentLocationID
	}

	if err := uc.exchangeUseCaseRepo.Create(&domain.Exchange{
		UserID:     userID,
		BookID:     bookID,
		Status:     domain.ExchangeRequested,
		ExpiresAt:  &expiresAt,
		LocationID: chosenLocationID,
	}); err != nil {
		return err
	}

	return nil
}

func (uc *BookUseCase) Borrow(bookID uuid.UUID, userID uuid.UUID) error {
	if err := validateBookID(bookID); err != nil {
		return err
	}
	if userID == uuid.Nil {
		return domain.NewValidationError("Неверный идентификатор пользователя")
	}

	book, err := uc.bookRepo.GetByID(bookID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.NewNotFoundError("Книга не найдена")
		}
		return err
	}
	if book.Status != domain.BookRequested {
		return domain.NewConflictError("Книга недоступна для выдачи")
	}
	exchanges, err := uc.exchangeUseCaseRepo.GetByBookID(bookID)
	if err != nil {
		return err
	}
	var requesterFound bool
	var matchedExchange *domain.Exchange
	for _, ex := range exchanges {
		if ex.UserID == userID && ex.Status == domain.ExchangeRequested {
			requesterFound = true
			matchedExchange = ex
			break
		}
	}
	if !requesterFound {
		return domain.NewForbiddenError("Только пользователь, который забронировал книгу, может её взять")
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
	// Update the related exchange status to 'borrowed' so it won't appear in reserved lists
	if matchedExchange != nil {
		matchedExchange.Status = domain.ExchangeBorrowed
		if err := uc.exchangeUseCaseRepo.Update(matchedExchange); err != nil {
			return err
		}
		// link movement to exchange and try to save movement update
		movement.ExchangeID = &matchedExchange.ID
		_ = uc.movementHistoryRepo.Update(movement)
	}
	return nil
}

func (uc *BookUseCase) Return(updatedBook *domain.Book, userID uuid.UUID, isAdmin bool) error {
	if updatedBook == nil {
		return domain.NewValidationError("Книга не передана")
	}
	if err := validateBookID(updatedBook.ID); err != nil {
		return err
	}
	if userID == uuid.Nil {
		return domain.NewValidationError("Неверный идентификатор пользователя")
	}

	bookFromDB, err := uc.bookRepo.GetByID(updatedBook.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.NewNotFoundError("Книга не найдена")
		}
		return err
	}

	if bookFromDB.Status != domain.BookBorrowed {
		return domain.NewConflictError("Только взятые книги могут быть возвращены")
	}

	exchanges, err := uc.exchangeUseCaseRepo.GetByBookID(bookFromDB.ID)
	if err != nil {
		return err
	}

	var borrowedExchange *domain.Exchange
	for _, ex := range exchanges {
		if ex == nil {
			continue
		}
		if ex.Status == domain.ExchangeBorrowed {
			borrowedExchange = ex
			break
		}
	}

	if borrowedExchange == nil {
		return domain.NewConflictError("Для книги не найдено активное заимствование")
	}

	if !isAdmin && borrowedExchange.UserID != userID {
		return domain.NewForbiddenError("Возвращать книгу может только текущий заёмщик или администратор")
	}

	title, err := validateTitle(updatedBook.Title)
	if err != nil {
		return err
	}
	author, err := validateAuthor(updatedBook.Author)
	if err != nil {
		return err
	}
	description, err := validateDescription(updatedBook.Description)
	if err != nil {
		return err
	}
	if err := validateCondition(updatedBook.Condition); err != nil {
		return err
	}
	if err := uc.validateLocation(updatedBook.CurrentLocationID); err != nil {
		return err
	}

	// Обновляем только нужные поля у объекта, который мы получили из БД
	bookFromDB.Status = domain.BookAvailable
	bookFromDB.Title = title
	bookFromDB.Author = author
	bookFromDB.Description = description
	if updatedBook.ImageURL != "" {
		imageURL, err := validateImageURL(updatedBook.ImageURL)
		if err != nil {
			return err
		}
		bookFromDB.ImageURL = imageURL
	}
	bookFromDB.CurrentLocationID = updatedBook.CurrentLocationID
	bookFromDB.Condition = updatedBook.Condition

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
	// Update any exchanges related to this book instead of deleting them.
	// Deleting would break FK constraints because movement history rows may reference exchanges.
	// For borrowed exchanges mark as returned and link the movement to that exchange.
	for _, ex := range exchanges {
		if ex == nil {
			continue
		}
		if ex.Status == domain.ExchangeBorrowed {
			ex.Status = domain.ExchangeReturned
			if uerr := uc.exchangeUseCaseRepo.Update(ex); uerr != nil {
				return uerr
			}
			movement.ExchangeID = &ex.ID
			if merr := uc.movementHistoryRepo.Update(movement); merr != nil {
				return merr
			}
		} else {
			ex.Status = domain.ExchangeCancelled
			if uerr := uc.exchangeUseCaseRepo.Update(ex); uerr != nil {
				return uerr
			}
		}
	}

	return nil
}

func (uc *BookUseCase) DeleteBook(bookID, userID uuid.UUID, isAdmin bool, reason string) error {
	if err := validateBookID(bookID); err != nil {
		return err
	}
	if userID == uuid.Nil {
		return domain.NewValidationError("Неверный идентификатор пользователя")
	}

	book, err := uc.bookRepo.GetByID(bookID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.NewNotFoundError("Книга не найдена")
		}
		return err
	}

	// Authorization: admin can delete any book. Non-admin can delete only if they
	// currently have the book borrowed (i.e. it's on their "shelf").
	if !isAdmin {
		// Check exchanges for this book to verify borrower
		exchanges, err := uc.exchangeUseCaseRepo.GetByBookID(bookID)
		if err != nil {
			return errors.New("Не удалось проверить статус заемщика")
		}
		var borrowerFound bool
		for _, ex := range exchanges {
			if ex == nil {
				continue
			}
			if ex.UserID == userID && ex.Status == domain.ExchangeBorrowed {
				borrowerFound = true
				break
			}
		}
		if !borrowerFound {
			return domain.NewForbiddenError("Доступ запрещён: удалять может только администратор или пользователь, у которого книга на полке")
		}
	}

	previousStatus := book.Status

	book.Status = domain.BookDeleted

	if err := uc.bookRepo.Update(book); err != nil {
		return err
	}

	movement := &domain.BookMovementHistory{
		BookID:         bookID,
		UserID:         &userID,
		Action:         "deleted",
		PreviousStatus: previousStatus,
		NewStatus:      domain.BookDeleted,
		Notes:          reason,
	}
	if err := uc.movementHistoryRepo.Create(movement); err != nil {
		return err
	}

	// Mark related exchanges as cancelled to preserve history and avoid FK errors
	if exchanges, err := uc.exchangeUseCaseRepo.GetByBookID(bookID); err == nil {
		for _, ex := range exchanges {
			if ex == nil {
				continue
			}
			ex.Status = domain.ExchangeCancelled
			if uerr := uc.exchangeUseCaseRepo.Update(ex); uerr != nil {
				return uerr
			}
		}
	}

	return nil

}

func (uc *BookUseCase) GetSummaryBooksList(userID *uuid.UUID) ([]*domain.BookSummary, error) {
	return uc.bookRepo.GetSummaryList(100, 0, userID)
}

func (uc *BookUseCase) GetBooksList(filter domain.BookListFilter) (*domain.BookListResponse, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}
	if filter.Order == "" {
		filter.Order = "desc"
	}
	if filter.OnlyAvailable {
		status := domain.BookAvailable
		filter.Status = &status
	}

	return uc.bookRepo.ListFiltered(filter)
}

func (uc *BookUseCase) UpdateBook(bookID uuid.UUID, userID uuid.UUID, isAdmin bool, input domain.BookUpdateInput) (*domain.Book, error) {
	if err := validateBookID(bookID); err != nil {
		return nil, err
	}
	if userID == uuid.Nil {
		return nil, domain.NewValidationError("Неверный идентификатор пользователя")
	}

	book, err := uc.bookRepo.GetByID(bookID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.NewNotFoundError("Книга не найдена")
		}
		return nil, err
	}

	allowed, err := uc.canManageBookMetadata(book, userID, isAdmin)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, domain.NewForbiddenError("Доступ запрещён: редактировать книгу может только владелец, текущий заёмщик или администратор")
	}

	if book.Status == domain.BookDeleted {
		return nil, domain.NewConflictError("Нельзя редактировать удалённую книгу")
	}

	if input.Title != nil {
		title, err := validateTitle(*input.Title)
		if err != nil {
			return nil, err
		}
		book.Title = title
	}

	if input.Author != nil {
		author, err := validateAuthor(*input.Author)
		if err != nil {
			return nil, err
		}
		book.Author = author
	}

	if input.Description != nil {
		description, err := validateDescription(*input.Description)
		if err != nil {
			return nil, err
		}
		book.Description = description
	}

	if input.Condition != nil {
		if err := validateCondition(*input.Condition); err != nil {
			return nil, err
		}
		book.Condition = *input.Condition
	}

	if input.ImageURL != nil {
		imageURL, err := validateImageURL(*input.ImageURL)
		if err != nil {
			return nil, err
		}
		book.ImageURL = imageURL
	}

	previousLocationID := book.CurrentLocationID
	if input.CurrentLocationIDSet {
		if err := uc.validateLocation(input.CurrentLocationID); err != nil {
			return nil, err
		}
		book.CurrentLocationID = input.CurrentLocationID
	}

	if err := uc.bookRepo.Update(book); err != nil {
		return nil, err
	}

	if input.CurrentLocationIDSet && !sameUUIDPtr(previousLocationID, book.CurrentLocationID) {
		movement := &domain.BookMovementHistory{
			BookID:         book.ID,
			FromLocationID: previousLocationID,
			ToLocationID:   book.CurrentLocationID,
			UserID:         &userID,
			Action:         "moved",
			Notes:          "Book location updated",
			PreviousStatus: book.Status,
			NewStatus:      book.Status,
		}
		if err := uc.movementHistoryRepo.Create(movement); err != nil {
			return nil, err
		}
	}

	return book, nil
}

func (uc *BookUseCase) GetBookByID(bookID uuid.UUID) (*domain.Book, error) {
	if err := validateBookID(bookID); err != nil {
		return nil, err
	}
	book, err := uc.bookRepo.GetByID(bookID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.NewNotFoundError("Книга не найдена")
		}
		return nil, err
	}
	return book, nil
}

func (uc *BookUseCase) canManageBookMetadata(book *domain.Book, userID uuid.UUID, isAdmin bool) (bool, error) {
	if book == nil {
		return false, domain.NewValidationError("Книга не передана")
	}
	if isAdmin || book.OwnerID == userID {
		return true, nil
	}
	if book.Status != domain.BookBorrowed {
		return false, nil
	}

	exchanges, err := uc.exchangeUseCaseRepo.GetByBookID(book.ID)
	if err != nil {
		return false, err
	}
	for _, ex := range exchanges {
		if ex == nil {
			continue
		}
		if ex.Status == domain.ExchangeBorrowed && ex.UserID == userID {
			return true, nil
		}
	}

	return false, nil
}

func sameUUIDPtr(a, b *uuid.UUID) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func (uc *BookUseCase) GetBookMovementHistory(bookID uuid.UUID) ([]*domain.BookMovementHistory, error) {
	if err := validateBookID(bookID); err != nil {
		return nil, err
	}
	_, err := uc.bookRepo.GetByID(bookID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.NewNotFoundError("Книга не найдена")
		}
		return nil, err
	}
	return uc.movementHistoryRepo.GetByBookID(bookID)
}

// GetBooksByOwner returns books owned by the provided user
func (uc *BookUseCase) GetBooksByOwner(ownerID uuid.UUID) ([]*domain.Book, error) {
	return uc.bookRepo.GetByOwner(ownerID)
}

// GetBooksFromExchangesByStatus returns books related to exchanges for user filtered by exchange status
func (uc *BookUseCase) GetBooksFromExchangesByStatus(userID uuid.UUID, status domain.ExchangeStatus) ([]*domain.Book, error) {
	exchanges, err := uc.exchangeUseCaseRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	var books []*domain.Book
	for _, ex := range exchanges {
		if ex == nil {
			continue
		}
		if ex.Status == status {
			// Exchange.Book is a struct (not a pointer) in domain; take its address
			books = append(books, &ex.Book)
		}
	}
	return books, nil
}

// GetExchangesFromUserByStatus returns exchanges for a user filtered by status
func (uc *BookUseCase) GetExchangesFromUserByStatus(userID uuid.UUID, status domain.ExchangeStatus) ([]*domain.Exchange, error) {
	exchanges, err := uc.exchangeUseCaseRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	var out []*domain.Exchange
	for _, ex := range exchanges {
		if ex == nil {
			continue
		}
		if ex.Status == status {
			out = append(out, ex)
		}
	}
	return out, nil
}

// GetBookStats returns aggregated statistics for a single book
func (uc *BookUseCase) GetBookStats(bookID uuid.UUID, userID uuid.UUID, isAdmin bool) (*domain.MyBooksStats, error) {
	if err := validateBookID(bookID); err != nil {
		return nil, err
	}
	if userID == uuid.Nil {
		return nil, domain.NewValidationError("Неверный идентификатор пользователя")
	}
	book, err := uc.bookRepo.GetByID(bookID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.NewNotFoundError("Книга не найдена")
		}
		return nil, err
	}
	if !isAdmin && book.OwnerID != userID {
		return nil, domain.NewForbiddenError("Доступ запрещён: статистика доступна только владельцу книги или администратору")
	}

	stats := &domain.MyBooksStats{
		TotalBooks:      1,
		StatusCounts:    map[string]int64{},
		ConditionCounts: map[string]int64{},
	}

	stats.StatusCounts[string(book.Status)] = 1
	stats.ConditionCounts[string(book.Condition)] = 1
	// set current status/condition for this book
	s := string(book.Status)
	c := string(book.Condition)
	stats.CurrentStatus = &s
	stats.CurrentCondition = &c

	// unique borrowers
	exchanges, err := uc.exchangeUseCaseRepo.GetByBookID(bookID)
	if err == nil {
		set := map[string]struct{}{}
		for _, ex := range exchanges {
			if ex == nil {
				continue
			}
			set[ex.UserID.String()] = struct{}{}
		}
		stats.TotalUniqueBorrowers = int64(len(set))
	}

	// movement history -> first/last location
	history, err := uc.movementHistoryRepo.GetByBookID(bookID)
	if err == nil && len(history) > 0 {
		var earliest *time.Time
		var latest *time.Time
		var earliestLoc *string
		var latestLoc *string
		for _, h := range history {
			if h == nil {
				continue
			}
			t := h.CreatedAt
			if earliest == nil || t.Before(*earliest) {
				earliest = &t
				if h.FromLocation != nil && h.FromLocation.Name != "" {
					name := h.FromLocation.Name
					earliestLoc = &name
				} else if h.ToLocation != nil && h.ToLocation.Name != "" {
					name := h.ToLocation.Name
					earliestLoc = &name
				}
			}
			if latest == nil || t.After(*latest) {
				latest = &t
				if h.ToLocation != nil && h.ToLocation.Name != "" {
					name := h.ToLocation.Name
					latestLoc = &name
				} else if h.FromLocation != nil && h.FromLocation.Name != "" {
					name := h.FromLocation.Name
					latestLoc = &name
				}
			}
		}
		if earliestLoc != nil {
			stats.FirstLocation = earliestLoc
		}
		if latestLoc != nil {
			stats.LastLocation = latestLoc
		}
	}

	return stats, nil
}

// SetBookImage sets the ImageURL for a given book.
func (uc *BookUseCase) SetBookImage(bookID uuid.UUID, userID uuid.UUID, isAdmin bool, imageURL string) error {
	if err := validateBookID(bookID); err != nil {
		return err
	}
	if userID == uuid.Nil {
		return domain.NewValidationError("Неверный идентификатор пользователя")
	}
	validatedImageURL, err := validateImageURL(imageURL)
	if err != nil {
		return err
	}

	book, err := uc.bookRepo.GetByID(bookID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.NewNotFoundError("Книга не найдена")
		}
		return err
	}
	allowed, err := uc.canManageBookMetadata(book, userID, isAdmin)
	if err != nil {
		return err
	}
	if !allowed {
		return domain.NewForbiddenError("Доступ запрещён: изменять изображение книги может только владелец, текущий заёмщик или администратор")
	}
	if strings.TrimSpace(book.ImageURL) != "" && book.ImageURL != validatedImageURL {
		if oldImagePath, pathErr := imagePathFromURL(book.ImageURL); pathErr == nil {
			if removeErr := os.Remove(oldImagePath); removeErr != nil && !os.IsNotExist(removeErr) {
				return removeErr
			}
		}
	}
	book.ImageURL = validatedImageURL
	if err := uc.bookRepo.Update(book); err != nil {
		return err
	}
	return nil
}

func (uc *BookUseCase) DeleteBookImage(bookID uuid.UUID, userID uuid.UUID, isAdmin bool) error {
	if err := validateBookID(bookID); err != nil {
		return err
	}
	if userID == uuid.Nil {
		return domain.NewValidationError("Неверный идентификатор пользователя")
	}

	book, err := uc.bookRepo.GetByID(bookID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.NewNotFoundError("Книга не найдена")
		}
		return err
	}
	allowed, err := uc.canManageBookMetadata(book, userID, isAdmin)
	if err != nil {
		return err
	}
	if !allowed {
		return domain.NewForbiddenError("Доступ запрещён: удалять изображение книги может только владелец, текущий заёмщик или администратор")
	}
	if strings.TrimSpace(book.ImageURL) == "" {
		return domain.NewValidationError("У книги нет изображения")
	}

	imagePath, err := imagePathFromURL(book.ImageURL)
	if err != nil {
		return err
	}
	if err := os.Remove(imagePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	book.ImageURL = ""
	if err := uc.bookRepo.Update(book); err != nil {
		return err
	}

	return nil
}

// GetMyBooksStats returns aggregated statistics for the given owner
func (uc *BookUseCase) GetMyBooksStats(ownerID uuid.UUID) (*domain.MyBooksStats, error) {
	books, err := uc.bookRepo.GetByOwner(ownerID)
	if err != nil {
		return nil, err
	}

	stats := &domain.MyBooksStats{
		TotalBooks:      int64(len(books)),
		StatusCounts:    map[string]int64{},
		ConditionCounts: map[string]int64{},
	}

	// set for unique borrowers
	borrowerSet := map[string]struct{}{}

	var earliestTime *time.Time
	var latestTime *time.Time
	var earliestLoc *string
	var latestLoc *string

	for _, b := range books {
		// counts
		stats.StatusCounts[string(b.Status)]++
		stats.ConditionCounts[string(b.Condition)]++

		// exchanges -> unique borrowers
		exchanges, err := uc.exchangeUseCaseRepo.GetByBookID(b.ID)
		if err == nil {
			for _, ex := range exchanges {
				if ex == nil {
					continue
				}
				borrowerSet[ex.UserID.String()] = struct{}{}
			}
		}

		// movement history -> first/last location
		history, err := uc.movementHistoryRepo.GetByBookID(b.ID)
		if err == nil && len(history) > 0 {
			// history is expected to be ordered by CreatedAt desc (most recent first) elsewhere,
			// but to be safe, iterate and track min/max
			for _, h := range history {
				if h == nil {
					continue
				}
				t := h.CreatedAt
				if earliestTime == nil || t.Before(*earliestTime) {
					earliestTime = &t
					if h.FromLocation != nil && h.FromLocation.Name != "" {
						name := h.FromLocation.Name
						earliestLoc = &name
					} else if h.ToLocation != nil && h.ToLocation.Name != "" {
						name := h.ToLocation.Name
						earliestLoc = &name
					}
				}
				if latestTime == nil || t.After(*latestTime) {
					latestTime = &t
					if h.ToLocation != nil && h.ToLocation.Name != "" {
						name := h.ToLocation.Name
						latestLoc = &name
					} else if h.FromLocation != nil && h.FromLocation.Name != "" {
						name := h.FromLocation.Name
						latestLoc = &name
					}
				}
			}
		}
	}

	stats.TotalUniqueBorrowers = int64(len(borrowerSet))
	if earliestLoc != nil {
		stats.FirstLocation = earliestLoc
	}
	if latestLoc != nil {
		stats.LastLocation = latestLoc
	}

	return stats, nil
}

// ExtendReservation attempts to extend an active reservation by 24 hours with constraints
func (uc *BookUseCase) ExtendReservation(bookID uuid.UUID, userID uuid.UUID) error {
	if err := validateBookID(bookID); err != nil {
		return err
	}
	if userID == uuid.Nil {
		return domain.NewValidationError("Неверный идентификатор пользователя")
	}

	exchanges, err := uc.exchangeUseCaseRepo.GetByBookID(bookID)
	if err != nil {
		return err
	}
	var target *domain.Exchange
	for _, ex := range exchanges {
		if ex == nil {
			continue
		}
		if ex.UserID == userID && ex.Status == domain.ExchangeRequested {
			target = ex
			break
		}
	}
	if target == nil {
		return domain.NewNotFoundError("Активное бронирование не найдено")
	}
	if target.ExpiresAt == nil {
		return domain.NewConflictError("У бронирования не задан срок действия")
	}
	if time.Now().After(*target.ExpiresAt) {
		return domain.NewConflictError("Срок бронирования уже истёк")
	}

	// Constraint: total reservation period cannot exceed 7 days
	maxDuration := 7 * 24 * time.Hour
	if target.ExpiresAt.Sub(target.BookedAt) >= maxDuration {
		return domain.NewConflictError("Нельзя продлить бронирование сверх максимального срока")
	}

	// extend by 24 hours
	newExp := target.ExpiresAt.Add(24 * time.Hour)
	// ensure we don't exceed max
	if newExp.Sub(target.BookedAt) > maxDuration {
		newExp = target.BookedAt.Add(maxDuration)
	}
	target.ExpiresAt = &newExp

	if err := uc.exchangeUseCaseRepo.Update(target); err != nil {
		return err
	}
	return nil
}

// CancelReservation cancels an active reservation and sets book status back to available
func (uc *BookUseCase) CancelReservation(bookID uuid.UUID, userID uuid.UUID) error {
	if err := validateBookID(bookID); err != nil {
		return err
	}
	if userID == uuid.Nil {
		return domain.NewValidationError("Неверный идентификатор пользователя")
	}

	exchanges, err := uc.exchangeUseCaseRepo.GetByBookID(bookID)
	if err != nil {
		return err
	}
	var target *domain.Exchange
	for _, ex := range exchanges {
		if ex == nil {
			continue
		}
		if ex.UserID == userID && ex.Status == domain.ExchangeRequested {
			target = ex
			break
		}
	}
	if target == nil {
		return domain.NewNotFoundError("Активное бронирование не найдено")
	}

	// mark exchange cancelled
	target.Status = domain.ExchangeCancelled
	if err := uc.exchangeUseCaseRepo.Update(target); err != nil {
		return err
	}

	// set book status back to available
	book, err := uc.bookRepo.GetByID(bookID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.NewNotFoundError("Книга не найдена")
		}
		return err
	}
	previousStatus := book.Status
	book.Status = domain.BookAvailable
	if err := uc.bookRepo.Update(book); err != nil {
		return err
	}

	// add movement history entry
	movement := &domain.BookMovementHistory{
		BookID:         bookID,
		UserID:         &userID,
		Action:         "reservation_cancelled",
		PreviousStatus: previousStatus,
		NewStatus:      domain.BookAvailable,
		Notes:          "Reservation cancelled by user",
	}
	if err := uc.movementHistoryRepo.Create(movement); err != nil {
		return err
	}

	return nil
}
