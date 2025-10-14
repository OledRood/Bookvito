package usecase
package usecase

import (
	"bookvito/internal/domain"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo domain.UserRepository
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(userRepo domain.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

// RegisterUser creates a new user with hashed password
func (uc *UserUseCase) RegisterUser(email, username, password, fullName, location string) (*domain.User, error) {
	// Check if user already exists
	existingUser, _ := uc.userRepo.GetByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	existingUser, _ = uc.userRepo.GetByUsername(username)
	if existingUser != nil {
		return nil, errors.New("user with this username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:    email,
		Username: username,
		Password: string(hashedPassword),
		FullName: fullName,
		Location: location,
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// LoginUser authenticates a user
func (uc *UserUseCase) LoginUser(email, password string) (*domain.User, error) {
	user, err := uc.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (uc *UserUseCase) GetUserByID(id uint) (*domain.User, error) {
	return uc.userRepo.GetByID(id)
}

// UpdateUser updates user information
func (uc *UserUseCase) UpdateUser(user *domain.User) error {
	return uc.userRepo.Update(user)
}

// DeleteUser deletes a user
func (uc *UserUseCase) DeleteUser(id uint) error {
	return uc.userRepo.Delete(id)
}

// ListUsers retrieves a list of users
func (uc *UserUseCase) ListUsers(limit, offset int) ([]*domain.User, error) {
	return uc.userRepo.List(limit, offset)
}
