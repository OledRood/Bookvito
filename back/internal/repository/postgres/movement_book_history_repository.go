package postgres

import (
	"bookvito/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type bookMovementHistoryRepository struct {
	db *gorm.DB
}

// NewBookMovementHistoryRepository creates a new book movement history repository
func NewBookMovementHistoryRepository(db *gorm.DB) domain.BookMovementHistoryRepository {
	return &bookMovementHistoryRepository{db: db}
}

// Create creates a new book movement history record
func (r *bookMovementHistoryRepository) Create(movement *domain.BookMovementHistory) error {
	return r.db.Create(movement).Error
}

// Update updates an existing book movement history record
func (r *bookMovementHistoryRepository) Update(movement *domain.BookMovementHistory) error {
	return r.db.Save(movement).Error
}

// GetByID retrieves a book movement history record by ID
func (r *bookMovementHistoryRepository) GetByID(id uuid.UUID) (*domain.BookMovementHistory, error) {
	var movement domain.BookMovementHistory
	err := r.db.
		Preload("Book").
		Preload("FromLocation").
		Preload("ToLocation").
		Preload("Exchange").
		Preload("User").
		First(&movement, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &movement, nil
}

// GetByBookID retrieves all movement history for a specific book
func (r *bookMovementHistoryRepository) GetByBookID(bookID uuid.UUID) ([]*domain.BookMovementHistory, error) {
	var movements []*domain.BookMovementHistory
	err := r.db.
		Where("book_id = ?", bookID).
		Preload("FromLocation").
		Preload("ToLocation").
		Preload("Exchange").
		Preload("User").
		Order("created_at DESC").
		Find(&movements).Error
	if err != nil {
		return nil, err
	}
	return movements, nil
}

// GetByExchangeID retrieves all movement history for a specific exchange
func (r *bookMovementHistoryRepository) GetByExchangeID(exchangeID uuid.UUID) ([]*domain.BookMovementHistory, error) {
	var movements []*domain.BookMovementHistory
	err := r.db.
		Where("exchange_id = ?", exchangeID).
		Preload("Book").
		Preload("FromLocation").
		Preload("ToLocation").
		Preload("User").
		Order("created_at DESC").
		Find(&movements).Error
	if err != nil {
		return nil, err
	}
	return movements, nil
}

// GetByUserID retrieves all movement history initiated by a specific user
func (r *bookMovementHistoryRepository) GetByUserID(userID uuid.UUID) ([]*domain.BookMovementHistory, error) {
	var movements []*domain.BookMovementHistory
	err := r.db.
		Where("user_id = ?", userID).
		Preload("Book").
		Preload("FromLocation").
		Preload("ToLocation").
		Preload("Exchange").
		Order("created_at DESC").
		Find(&movements).Error
	if err != nil {
		return nil, err
	}
	return movements, nil
}

// List retrieves a paginated list of movement history
func (r *bookMovementHistoryRepository) List(limit, offset int) ([]*domain.BookMovementHistory, error) {
	var movements []*domain.BookMovementHistory
	err := r.db.
		Preload("Book").
		Preload("FromLocation").
		Preload("ToLocation").
		Preload("Exchange").
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&movements).Error
	if err != nil {
		return nil, err
	}
	return movements, nil
}
