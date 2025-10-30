package usecase

import (
	"bookvito/internal/domain"
	"errors"
	"fmt"
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

func (uc *BookUseCase) Request(bookID uuid.UUID, userID uuid.UUID, locationID *uuid.UUID) error {
	book, err := uc.bookRepo.GetByID(bookID)
	if err != nil {
		return err
	}
	if book.Status != domain.BookAvailable || book.Status == "" {
		return fmt.Errorf("Книга недоступна для бронирования (статус=%s)", book.Status)
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
	book, err := uc.bookRepo.GetByID(bookID)
	if err != nil {
		return err
	}
	if book.Status != domain.BookRequested {
		return errors.New("Книга недоступна для выдачи")
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
		return errors.New("Только пользователь, который забронировал книгу, может её взять")
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

func (uc *BookUseCase) Return(updatedBook *domain.Book, userID uuid.UUID) error {
	if updatedBook.Title == "" {
		return errors.New("Название книги не может быть пустым")
	}
	if updatedBook.Author == "" {
		return errors.New("Автор книги не может быть пустым")
	}
	println("Returning book status:", updatedBook.Status)

	bookFromDB, err := uc.bookRepo.GetByID(updatedBook.ID)
	if err != nil {
		return errors.New("Книга не найдена в базе")
	}

	if bookFromDB.Status != domain.BookBorrowed {
		return errors.New("Только взятые книги могут быть возвращены")
	}

	// Обновляем только нужные поля у объекта, который мы получили из БД
	bookFromDB.Status = domain.BookAvailable
	bookFromDB.Title = updatedBook.Title
	bookFromDB.Author = updatedBook.Author
	bookFromDB.Description = updatedBook.Description
	// Preserve existing ImageURL if the incoming payload doesn't include one.
	// Some clients don't send image_url when returning a book, and unconditional
	// assignment would clear the stored image. Update only when a non-empty
	// value is provided.
	if updatedBook.ImageURL != "" {
		bookFromDB.ImageURL = updatedBook.ImageURL
	}
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
	// Update any exchanges related to this book instead of deleting them.
	// Deleting would break FK constraints because movement history rows may reference exchanges.
	// For borrowed exchanges mark as returned and link the movement to that exchange.
	if exchanges, err := uc.exchangeUseCaseRepo.GetByBookID(bookFromDB.ID); err == nil {
		for _, ex := range exchanges {
			if ex == nil {
				continue
			}
			if ex.Status == domain.ExchangeBorrowed {
				ex.Status = domain.ExchangeReturned
				if uerr := uc.exchangeUseCaseRepo.Update(ex); uerr != nil {
					return uerr
				}
				// link the return movement to the exchange
				movement.ExchangeID = &ex.ID
				// best-effort update; if it fails we still continue but return the error
				if merr := uc.movementHistoryRepo.Update(movement); merr != nil {
					return merr
				}
			} else {
				// For other statuses (e.g. requested), mark as cancelled to preserve history
				ex.Status = domain.ExchangeCancelled
				if uerr := uc.exchangeUseCaseRepo.Update(ex); uerr != nil {
					return uerr
				}
			}
		}
	}

	return nil
}

func (uc *BookUseCase) DeleteBook(bookID, userID uuid.UUID, isAdmin bool, reason string) error {
	book, err := uc.bookRepo.GetByID(bookID)

	if err != nil {
		return errors.New("Книга не найдена")
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
			return errors.New("Доступ запрещён: удалять может только администратор или пользователь, у которого книга на полке")
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

func (uc *BookUseCase) GetBooksList() ([]*domain.Book, error) {
	return uc.bookRepo.List(100, 0)

}

// SearchBooks performs a text search using repository.Search and then applies
// the same visibility filtering as GetSummaryBooksList: only available books,
// exclude books owned by the requesting user and exclude books the user has
// returned before. Returns BookSummary objects for the UI.
func (uc *BookUseCase) SearchBooks(query string, limit, offset int, userID *uuid.UUID) ([]*domain.BookSummary, error) {
	if limit == 0 {
		limit = 100
	}

	books, err := uc.bookRepo.Search(query, limit, offset)
	if err != nil {
		return nil, err
	}

	// Build set of returned book IDs for user (if any)
	returned := map[uuid.UUID]struct{}{}
	if userID != nil {
		exchanges, err := uc.exchangeUseCaseRepo.GetByUserID(*userID)
		if err == nil {
			for _, ex := range exchanges {
				if ex == nil {
					continue
				}
				if ex.Status == domain.ExchangeReturned {
					returned[ex.BookID] = struct{}{}
				}
			}
		}
	}

	var out []*domain.BookSummary
	for _, b := range books {
		if b == nil {
			continue
		}
		// Follow summary visibility rules: only available books
		if b.Status != domain.BookAvailable {
			continue
		}
		// exclude owner's own books when user is provided
		if userID != nil && b.OwnerID == *userID {
			continue
		}
		// exclude returned books for user
		if _, ok := returned[b.ID]; ok {
			continue
		}

		out = append(out, &domain.BookSummary{
			ID:       b.ID,
			ImageURL: b.ImageURL,
			Title:    b.Title,
			Author:   b.Author,
		})
	}

	return out, nil
}

func (uc *BookUseCase) GetBookByID(bookID uuid.UUID) (*domain.Book, error) {
	return uc.bookRepo.GetByID(bookID)
}

func (uc *BookUseCase) GetBookMovementHistory(bookID uuid.UUID) ([]*domain.BookMovementHistory, error) {
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
func (uc *BookUseCase) GetBookStats(bookID uuid.UUID) (*domain.MyBooksStats, error) {
	book, err := uc.bookRepo.GetByID(bookID)
	if err != nil {
		return nil, err
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
func (uc *BookUseCase) SetBookImage(bookID uuid.UUID, imageURL string) error {
	book, err := uc.bookRepo.GetByID(bookID)
	if err != nil {
		return err
	}
	book.ImageURL = imageURL
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
		return errors.New("no active reservation found for this user and book")
	}
	if target.ExpiresAt == nil {
		return errors.New("reservation has no expiry set")
	}
	if time.Now().After(*target.ExpiresAt) {
		return errors.New("reservation already expired")
	}

	// Constraint: total reservation period cannot exceed 7 days
	maxDuration := 7 * 24 * time.Hour
	if target.ExpiresAt.Sub(target.BookedAt) >= maxDuration {
		return errors.New("maximum reservation duration exceeded; cannot extend")
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
		return errors.New("no active reservation found for this user and book")
	}

	// mark exchange cancelled
	target.Status = domain.ExchangeCancelled
	if err := uc.exchangeUseCaseRepo.Update(target); err != nil {
		return err
	}

	// set book status back to available
	book, err := uc.bookRepo.GetByID(bookID)
	if err != nil {
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
