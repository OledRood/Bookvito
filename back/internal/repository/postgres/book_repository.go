package postgres

import (
	"bookvito/internal/domain"

	"gorm.io/gorm"
)

type bookRepository struct {
	db *gorm.DB
}

// NewBookRepository creates a new book repository
func NewBookRepository(db *gorm.DB) domain.BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(book *domain.Book) error {
	return r.db.Create(book).Error
}

func (r *bookRepository) GetByID(id uint) (*domain.Book, error) {
	var book domain.Book
	err := r.db.Preload("Owner").First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) GetByOwnerID(ownerID uint) ([]*domain.Book, error) {
	var books []*domain.Book
	err := r.db.Where("owner_id = ?", ownerID).Find(&books).Error
	return books, err
}

func (r *bookRepository) Update(book *domain.Book) error {
	return r.db.Save(book).Error
}

func (r *bookRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Book{}, id).Error
}

func (r *bookRepository) List(limit, offset int) ([]*domain.Book, error) {
	var books []*domain.Book
	err := r.db.Preload("Owner").Limit(limit).Offset(offset).Find(&books).Error
	return books, err
}

func (r *bookRepository) Search(query string, limit, offset int) ([]*domain.Book, error) {
	var books []*domain.Book
	searchPattern := "%" + query + "%"
	err := r.db.Preload("Owner").
		Where("title ILIKE ? OR author ILIKE ? OR description ILIKE ?", searchPattern, searchPattern, searchPattern).
		Limit(limit).
		Offset(offset).
		Find(&books).Error
	return books, err
}

func (r *bookRepository) GetAvailable(limit, offset int) ([]*domain.Book, error) {
	var books []*domain.Book
	err := r.db.Preload("Owner").
		Where("available = ?", true).
		Limit(limit).
		Offset(offset).
		Find(&books).Error
	return books, err
}
