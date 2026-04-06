package usecase

import (
	"bookvito/internal/domain"
	"testing"
)

func TestLocationUseCase_CreateAndGetByID(t *testing.T) {
	locationRepo := newFakeLocationRepo()
	uc := NewLocationUseCase(locationRepo)
	location := &domain.Location{Name: "Central", Address: "Main street"}

	if err := uc.Create(location); err != nil {
		t.Fatalf("create location: %v", err)
	}

	stored, err := uc.GetByID(location.ID)
	if err != nil {
		t.Fatalf("get location: %v", err)
	}
	if stored.Name != location.Name || stored.Address != location.Address {
		t.Fatalf("unexpected location: %+v", stored)
	}
}
