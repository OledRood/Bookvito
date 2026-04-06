package usecase

import (
	"bookvito/internal/domain"
	"testing"
)

func TestModerUseCase_CreateReport_EmptyReasonValidation(t *testing.T) {
	uc := NewModerUseCase(&fakeReportRepo{}, newFakeBookRepo())

	err := uc.CreateReport(mustUUID(), mustUUID(), "")
	if !hasErrorCode(err, domain.ErrorCodeValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestModerUseCase_ArchiveBook_NotFound(t *testing.T) {
	uc := NewModerUseCase(&fakeReportRepo{}, newFakeBookRepo())

	err := uc.ArchiveBook(mustUUID())
	if !hasErrorCode(err, domain.ErrorCodeNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
}
