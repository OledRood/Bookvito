package postgres

import (
	"bookvito/internal/domain"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type locationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) domain.LocationRepository {
	return &locationRepository{db: db}
}

func (r *locationRepository) Create(location *domain.Location) error {
	existing, err := r.GetByAddress(location.Address)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existing != nil {
		location.ID = existing.ID
	}
	return r.db.Create(location).Error
}

func (r *locationRepository) GetByID(id uuid.UUID) (*domain.Location, error) {
	var location domain.Location
	err := r.db.Preload("Books").First(&location, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &location, nil
}

func (r *locationRepository) GetByAddress(address string) (*domain.Location, error) {
	var location domain.Location
	err := r.db.Preload("Books").First(&location, "address = ?", address).Error
	if err != nil {
		return nil, err
	}
	return &location, nil
}

func (r *locationRepository) GetAll() ([]domain.Location, error) {
	var locations []domain.Location
	err := r.db.Preload("Books").Find(&locations).Error
	return locations, err
}

func (r *locationRepository) Update(location *domain.Location) error {
	return r.db.Save(location).Error
}

func (r *locationRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&domain.Location{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
