package usecase

import (
	"bookvito/internal/domain"
)

type AdminUseCase struct {
	userRepo     domain.UserRepository
	bookRepo     domain.BookRepository
	exchangeRepo domain.ExchangeRepository
	reportRepo   domain.ReportRepository
}

func NewAdminUseCase(
	userRepo domain.UserRepository,
	bookRepo domain.BookRepository,
	exchangeRepo domain.ExchangeRepository,
	reportRepo domain.ReportRepository,
) *AdminUseCase {
	return &AdminUseCase{
		userRepo:     userRepo,
		bookRepo:     bookRepo,
		exchangeRepo: exchangeRepo,
		reportRepo:   reportRepo,
	}
}

func (uc *AdminUseCase) GetStats() (*domain.AdminStats, error) {
	users, err := uc.userRepo.List(0, 0)
	if err != nil {
		return nil, err
	}

	books, err := uc.bookRepo.List(0, 0)
	if err != nil {
		return nil, err
	}

	exchanges, err := uc.exchangeRepo.List(0, 0)
	if err != nil {
		return nil, err
	}
	activeExchanges := int64(0)
	for _, e := range exchanges {
		if e.Status == domain.ExchangeRequested || e.Status == domain.ExchangeBorrowed {
			activeExchanges++
		}
	}

	pendingReports, err := uc.reportRepo.CountPending()
	if err != nil {
		return nil, err
	}

	return &domain.AdminStats{
		TotalUsers:      int64(len(users)),
		TotalBooks:      int64(len(books)),
		ActiveExchanges: activeExchanges,
		PendingReports:  pendingReports,
	}, nil
}
