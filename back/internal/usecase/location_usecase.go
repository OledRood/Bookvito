package usecase

import (
	"bookvito/internal/domain"

	"github.com/google/uuid"
	// "golang.org/x/crypto/bcrypt"
)

type LocationUseCase struct {
	locationRepo domain.LocationRepository
}

func NewLocationUseCase(locationRepo domain.LocationRepository) *LocationUseCase {
	return &LocationUseCase{
		locationRepo: locationRepo,
	}
}

func (uc *LocationUseCase) Create(location *domain.Location) error {
	return uc.locationRepo.Create(location)
}

func (uc *LocationUseCase) GetByID(id uuid.UUID) (*domain.Location, error) {
	return uc.locationRepo.GetByID(id)
}

func (uc *LocationUseCase) GetAll() ([]domain.Location, error) {
	return uc.locationRepo.GetAll()
}

func (uc *LocationUseCase) Update(location *domain.Location) error {
	return uc.locationRepo.Update(location)
}

func (uc *LocationUseCase) Delete(id uuid.UUID) error {
	return uc.locationRepo.Delete(id)
}
