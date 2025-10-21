package postgres

import (
	"bookvito/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type exchangeRepository struct {
	db *gorm.DB
}

// NewExchangeRepository creates a new exchange repository
func NewExchangeRepository(db *gorm.DB) domain.ExchangeRepository {
	return &exchangeRepository{db: db}
}

func (r *exchangeRepository) Create(exchange *domain.Exchange) error {
	return r.db.Create(exchange).Error
}

func (r *exchangeRepository) GetByID(id uuid.UUID) (*domain.Exchange, error) {
	var exchange domain.Exchange
	err := r.db.Preload("User").Preload("Book").Preload("Location").First(&exchange, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &exchange, nil
}

func (r *exchangeRepository) GetByUserID(userID uuid.UUID) ([]*domain.Exchange, error) {
	var exchanges []*domain.Exchange
	// Preload Book and Location, User is redundant as we are querying by user_id
	err := r.db.Preload("Book").Preload("Location").Where("user_id = ?", userID).Find(&exchanges).Error
	if err != nil {
		return nil, err
	}
	return exchanges, nil
}

func (r *exchangeRepository) GetByBookID(bookID uuid.UUID) ([]*domain.Exchange, error) {
	var exchanges []*domain.Exchange
	// Preload User and Location, Book is redundant as we are querying by book_id
	err := r.db.Preload("User").Preload("Location").Where("book_id = ?", bookID).Find(&exchanges).Error
	if err != nil {
		return nil, err
	}
	return exchanges, nil
}

func (r *exchangeRepository) Update(exchange *domain.Exchange) error {
	return r.db.Save(exchange).Error
}

func (r *exchangeRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Exchange{}, "id = ?", id).Error
}

func (r *exchangeRepository) List(limit, offset int) ([]*domain.Exchange, error) {
	var exchanges []*domain.Exchange
	err := r.db.Preload("User").Preload("Book").Preload("Location").Limit(limit).Offset(offset).Find(&exchanges).Error
	return exchanges, err
}
