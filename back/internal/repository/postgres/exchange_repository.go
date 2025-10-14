package postgres

import (
	"bookvito/internal/domain"

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

func (r *exchangeRepository) GetByID(id uint) (*domain.Exchange, error) {
	var exchange domain.Exchange
	err := r.db.Preload("Requester").Preload("Owner").Preload("Book").First(&exchange, id).Error
	if err != nil {
		return nil, err
	}
	return &exchange, nil
}

func (r *exchangeRepository) GetByRequesterID(requesterID uint) ([]*domain.Exchange, error) {
	var exchanges []*domain.Exchange
	err := r.db.Preload("Requester").Preload("Owner").Preload("Book").
		Where("requester_id = ?", requesterID).
		Find(&exchanges).Error
	return exchanges, err
}

func (r *exchangeRepository) GetByOwnerID(ownerID uint) ([]*domain.Exchange, error) {
	var exchanges []*domain.Exchange
	err := r.db.Preload("Requester").Preload("Owner").Preload("Book").
		Where("owner_id = ?", ownerID).
		Find(&exchanges).Error
	return exchanges, err
}

func (r *exchangeRepository) GetByBookID(bookID uint) ([]*domain.Exchange, error) {
	var exchanges []*domain.Exchange
	err := r.db.Preload("Requester").Preload("Owner").Preload("Book").
		Where("book_id = ?", bookID).
		Find(&exchanges).Error
	return exchanges, err
}

func (r *exchangeRepository) Update(exchange *domain.Exchange) error {
	return r.db.Save(exchange).Error
}

func (r *exchangeRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Exchange{}, id).Error
}

func (r *exchangeRepository) List(limit, offset int) ([]*domain.Exchange, error) {
	var exchanges []*domain.Exchange
	err := r.db.Preload("Requester").Preload("Owner").Preload("Book").
		Limit(limit).
		Offset(offset).
		Find(&exchanges).Error
	return exchanges, err
}
