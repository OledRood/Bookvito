package postgres

import (
	"bookvito/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.User{}, "id = ?", id).Error
}

func (r *userRepository) List(limit, offset int) ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}
func (r *userRepository) GetByRefreshToken(refreshToken string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("refresh_token = ?", refreshToken).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
