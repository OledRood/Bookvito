package usecase

import (
	"bookvito/internal/domain"
	"log"
)

type ExchangeUseCase struct {
	exchangeRepo domain.ExchangeRepository
	bookRepo     domain.BookRepository
	userRepo     domain.UserRepository
	movementRepo domain.BookMovementHistoryRepository
}

// NewExchangeUseCase creates a new exchange use case
func NewExchangeUseCase(exchangeRepo domain.ExchangeRepository, bookRepo domain.BookRepository, userRepo domain.UserRepository, movementRepo domain.BookMovementHistoryRepository) *ExchangeUseCase {
	return &ExchangeUseCase{
		exchangeRepo: exchangeRepo,
		bookRepo:     bookRepo,
		userRepo:     userRepo,
		movementRepo: movementRepo,
	}
}

// CancelExpiredExchanges находит и отменяет все просроченные бронирования.
func (uc *ExchangeUseCase) CancelExpiredExchanges() error {
	expiredExchanges, err := uc.exchangeRepo.GetExpired()
	if err != nil {
		return err
	}
	if len(expiredExchanges) == 0 {
		log.Println("No expired exchanges found.")
		return nil
	}
	log.Printf("Found %d expired exchanges to cancel.", len(expiredExchanges))

	for _, exchange := range expiredExchanges {
		// 1. Обновляем статус бронирования на "отменено"
		exchange.Status = domain.ExchangeCancelled
		if err := uc.exchangeRepo.Update(exchange); err != nil {
			// Логируем ошибку, но продолжаем, чтобы не остановить весь процесс
			log.Printf("failed to update expired exchange %s: %v", exchange.ID, err)
			continue
		}

		// 2. Возвращаем книге статус "доступна"
		book, err := uc.bookRepo.GetByID(exchange.BookID)
		if err != nil {
			log.Printf("failed to get book %s for expired exchange %s: %v", exchange.BookID, exchange.ID, err)
			continue
		}

		// Убедимся, что мы не меняем статус книги, которая уже была взята или возвращена
		if book.Status == domain.BookRequested {
			book.Status = domain.BookAvailable
			if err := uc.bookRepo.Update(book); err != nil {
				log.Printf("failed to update book status for book %s: %v", book.ID, err)
				continue
			}

			// 3. Создаем запись в истории перемещений
			movement := &domain.BookMovementHistory{
				BookID:         book.ID,
				Action:         "request_cancelled",
				Notes:          "Book request expired and was automatically cancelled",
				PreviousStatus: domain.BookRequested,
				NewStatus:      domain.BookAvailable,
			}
			uc.movementRepo.Create(movement) // Ошибку можно логировать, но она не критична для основного потока
		}
	}
	return nil
}
