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

func (r *bookRepository) ListFiltered(filter domain.BookListFilter) (*domain.BookListResponse, error) {
	var books []*domain.Book

	sortColumns := map[string]string{
		"title":               "title",
		"author":              "author",
		"status":              "status",
		"created_at":          "created_at",
		"updated_at":          "updated_at",
		"current_location_id": "current_location_id",
	}

	sortColumn, ok := sortColumns[filter.SortBy]
	if !ok {
		sortColumn = "created_at"
	}

	order := "desc"
	if filter.Order == "asc" {
		order = "asc"
	}

	query := r.db.Preload("CurrentLocation").Model(&domain.Book{})

	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("title ILIKE ? OR author ILIKE ?", searchPattern, searchPattern)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.LocationID != nil {
		query = query.Where("current_location_id = ?", *filter.LocationID)
	}

	if filter.ExcludeUserID != nil {
		query = query.Where("owner_id != ?", *filter.ExcludeUserID)

		sub := r.db.Model(&domain.Exchange{}).
			Select("book_id").
			Where("user_id = ? AND status = ?", *filter.ExcludeUserID, domain.ExchangeReturned)
		query = query.Where("id NOT IN (?)", sub)
	}

	err := query.
		Order(sortColumn + " " + order).
		Limit(filter.Limit + 1).
		Offset(filter.Offset).
		Find(&books).Error
	if err != nil {
		return nil, err
	}

	hasMore := len(books) > filter.Limit
	if hasMore {
		books = books[:filter.Limit]
	}

	return &domain.BookListResponse{
		Items:   books,
		HasMore: hasMore,
	}, nil
}

func (r *bookRepository) GetSummaryList(limit, offset int, excludeUserID *uuid.UUID) ([]*domain.BookSummary, error) {
	var summaries []*domain.BookSummary

	query := r.db.Model(&domain.Book{}).
		Select("id, image_url, title, author").
		Where("status = ?", domain.BookAvailable)

	if excludeUserID != nil {
		// exclude books owned by the user
		query = query.Where("owner_id != ?", *excludeUserID)

		// exclude books that the user has returned (i.e. had exchanges with status 'returned')
		// We only want to hide books that the user has already read/returned; other interactions
		// (requested/borrowed) do not affect the public summary because only available books
		// are selected above.
		sub := r.db.Model(&domain.Exchange{}).
			Select("book_id").
			Where("user_id = ? AND status = ?", *excludeUserID, domain.ExchangeReturned)
		query = query.Where("id NOT IN (?)", sub)
	}

	err := query.Limit(limit).
		Offset(offset).
		Find(&summaries).Error
	return summaries, err
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

func (r *bookRepository) GetByOwner(ownerID uuid.UUID) ([]*domain.Book, error) {
	var books []*domain.Book
	err := r.db.Preload("CurrentLocation").
		Where("owner_id = ?", ownerID).
		Find(&books).Error
	return books, err
}
