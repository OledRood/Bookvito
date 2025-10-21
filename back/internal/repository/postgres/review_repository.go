package postgres

import (
	"bookvito/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) domain.ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(review *domain.Review) error {
	return r.db.Create(review).Error
}

func (r *reviewRepository) GetByID(id uuid.UUID) (*domain.Review, error) {
	var review domain.Review
	err := r.db.Preload("Book").Preload("User").First(&review, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) GetByBookID(bookID uuid.UUID) ([]domain.Review, error) {
	var reviews []domain.Review
	err := r.db.Preload("User").Where("book_id = ?", bookID).Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) GetByUserID(userID uuid.UUID) ([]domain.Review, error) {
	var reviews []domain.Review
	err := r.db.Preload("Book").Where("user_id = ?", userID).Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) Update(review *domain.Review) error {
	return r.db.Save(review).Error
}

func (r *reviewRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Review{}, "id = ?", id).Error
}
