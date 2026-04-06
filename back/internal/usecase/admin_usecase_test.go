package usecase

import (
	"bookvito/internal/domain"
	"testing"
)

func TestAdminUseCase_GetStats_CountsActiveExchangesAndReports(t *testing.T) {
	userRepo := newFakeUserRepo()
	bookRepo := newFakeBookRepo()
	exchangeRepo := &fakeExchangeRepo{
		items: []*domain.Exchange{
			{ID: mustUUID(), Status: domain.ExchangeRequested},
			{ID: mustUUID(), Status: domain.ExchangeBorrowed},
			{ID: mustUUID(), Status: domain.ExchangeReturned},
		},
	}
	reportRepo := &fakeReportRepo{
		reports: []*domain.Report{
			{ID: mustUUID(), Status: domain.ReportPending},
			{ID: mustUUID(), Status: domain.ReportResolved},
		},
	}

	_ = userRepo.Create(&domain.User{Email: "one@example.com", Name: "One"})
	_ = userRepo.Create(&domain.User{Email: "two@example.com", Name: "Two"})
	_ = bookRepo.Create(&domain.Book{OwnerID: mustUUID(), Title: "A", Author: "B"})

	uc := NewAdminUseCase(userRepo, bookRepo, exchangeRepo, reportRepo)

	stats, err := uc.GetStats()
	if err != nil {
		t.Fatalf("get stats: %v", err)
	}
	if stats.TotalUsers != 2 || stats.TotalBooks != 1 || stats.ActiveExchanges != 2 || stats.PendingReports != 1 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
}
