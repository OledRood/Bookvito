package postgres

import (
	"bookvito/internal/domain"

	"github.com/google/uuid"
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

func (r *bookRepository) GetByID(id uuid.UUID) (*domain.Book, error) {
	var book domain.Book
	err := r.db.Preload("CurrentLocation").Preload("Reviews").First(&book, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) Update(book *domain.Book) error {
	return r.db.Save(book).Error
}

func (r *bookRepository) Delete(bookID uuid.UUID) error {
	// Удаляем книгу только если пользователь является владельцем
	var book domain.Book
	if err := r.db.First(&book, "id = ?", bookID).Error; err != nil {
		return err
	}
	return r.db.Delete(&domain.Book{}, "id = ?", bookID).Error
}

func (r *bookRepository) List(limit, offset int) ([]*domain.Book, error) {
	var books []*domain.Book

	err := r.db.Model(&domain.Book{}). // Указываем модель, но выбираем только нужные поля
						Select("id, image_url, title, author"). // Выбираем только нужные поля
						Limit(limit).
						Offset(offset).
						Find(&books).Error
	return books, err
}

func (r *bookRepository) GetSummaryList(limit, offset int) ([]*domain.BookSummary, error) {
	var summaries []*domain.BookSummary
	err := r.db.Model(&domain.Book{}). // Указываем модель, но выбираем только нужные поля
						Select("id, image_url, title, author"). // Выбираем только нужные поля
						Limit(limit).
						Offset(offset).
						Find(&summaries).Error
	return summaries, err
}

func (r *bookRepository) Search(query string, limit, offset int) ([]*domain.Book, error) {
	var books []*domain.Book
	searchPattern := "%" + query + "%"
	err := r.db.Preload("CurrentLocation").
		Where("title ILIKE ? OR author ILIKE ? OR description ILIKE ?", searchPattern, searchPattern, searchPattern).
		Limit(limit).
		Offset(offset).
		Find(&books).Error
	return books, err
}

func (r *bookRepository) GetByStatus(status domain.BookStatus, limit, offset int) ([]*domain.Book, error) {
	var books []*domain.Book
	err := r.db.Preload("CurrentLocation").
		Where("status = ?", status).
		Limit(limit).
		Offset(offset).
		Find(&books).Error
	return books, err
}

func (r *bookRepository) GetByLocationID(locationID uuid.UUID) ([]*domain.Book, error) {
	var books []*domain.Book
	err := r.db.Preload("CurrentLocation").
		Where("current_location_id = ?", locationID).
		Find(&books).Error
	return books, err
}
