package usecase

import (
	"bookvito/internal/domain"
	"errors"

	"github.com/google/uuid"
)

type ModerUseCase struct {
	reportRepo domain.ReportRepository
	bookRepo   domain.BookRepository
}

func NewModerUseCase(reportRepo domain.ReportRepository, bookRepo domain.BookRepository) *ModerUseCase {
	return &ModerUseCase{
		reportRepo: reportRepo,
		bookRepo:   bookRepo,
	}
}

func (uc *ModerUseCase) GetReports(limit, offset int) ([]*domain.Report, error) {
	return uc.reportRepo.GetPending(limit, offset)
}

func (uc *ModerUseCase) CreateReport(bookID uuid.UUID, userID uuid.UUID, reason string) error {
	if reason == "" {
		return errors.New("причина жалобы обязательна")
	}
	report := &domain.Report{
		BookID: bookID,
		UserID: userID,
		Reason: reason,
		Status: domain.ReportPending,
	}
	return uc.reportRepo.Create(report)
}

func (uc *ModerUseCase) ResolveReport(reportID uuid.UUID) error {
	return uc.reportRepo.UpdateStatus(reportID, domain.ReportResolved)
}

func (uc *ModerUseCase) DismissReport(reportID uuid.UUID) error {
	return uc.reportRepo.UpdateStatus(reportID, domain.ReportDismissed)
}

func (uc *ModerUseCase) ArchiveBook(bookID uuid.UUID) error {
	book, err := uc.bookRepo.GetByID(bookID)
	if err != nil {
		return err
	}
	book.Status = domain.BookArchived
	return uc.bookRepo.Update(book)
}
