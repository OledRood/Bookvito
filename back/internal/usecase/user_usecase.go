package usecase

import (
	"bookvito/internal/domain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	// "golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo     domain.UserRepository
	movementRepo domain.BookMovementHistoryRepository
	jwtSecret    string
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(userRepo domain.UserRepository, movementRepo domain.BookMovementHistoryRepository, jwtSecret string) *UserUseCase {
	return &UserUseCase{
		userRepo:     userRepo,
		movementRepo: movementRepo,
		jwtSecret:    jwtSecret,
	}
}

func (uc *UserUseCase) RegisterUser(email string, password string, name string) (*domain.TokenResponse, error) {
	existingUser, _ := uc.userRepo.GetByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}
	err = uc.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// Generate new pair of tokens
	return uc.generateTokenPair(user)
}

func (uc *UserUseCase) generateTokenPair(user *domain.User) (*domain.TokenResponse, error) {
	// Generate Access Token
	claims := jwt.MapClaims{
		"sub":    user.ID,
		"userId": user.ID.String(),
		"email":  user.Email,
		"name":   user.Name,
		"role":   user.Role,
		"exp":    time.Now().Add(time.Hour * 10).Unix(),
		"iat":    time.Now().Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := jwtToken.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return nil, err
	}

	refreshToken := uuid.New().String()
	refreshTokenExpiresAt := time.Now().Add(time.Hour * 24 * 120) // Expires in 120 days

	user.RefreshToken = refreshToken
	user.RefreshTokenExpiresAt = refreshTokenExpiresAt

	if err := uc.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &domain.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *UserUseCase) LoginUser(email, password string) (*domain.TokenResponse, error) {
	user, err := uc.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}
	return uc.generateTokenPair(user)
}

func (uc *UserUseCase) RefreshToken(refreshToken string) (*domain.TokenResponse, error) {
	user, err := uc.userRepo.GetByRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if time.Now().After(user.RefreshTokenExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	// Generate new pair of tokens
	return uc.generateTokenPair(user)
}

// // LoginUser authenticates a user
// func (uc *UserUseCase) LoginUser(email, password string) (*domain.User, error) {
// 	user, err := uc.userRepo.GetByEmail(email)
// 	if err != nil {
// 		return nil, errors.New("invalid email or password")
// 	}

// 	// Compare password
// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
// 		return nil, errors.New("invalid email or password")
// 	}

// 	return user, nil
// }

// // GetUserByID retrieves a user by ID
// func (uc *UserUseCase) GetUserByID(id uuid.UUID) (*domain.User, error) {
// 	return uc.userRepo.GetByID(id)
// }

// // UpdateUser updates user information
// func (uc *UserUseCase) UpdateUser(user *domain.User) error {
// 	return uc.userRepo.Update(user)
// }

func (uc *UserUseCase) GetUserByID(id string) (*domain.User, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}
	return uc.userRepo.GetByID(uuidID)
}

func (uc *UserUseCase) DeleteUser(id string) error {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid user ID format")
	}
	return uc.userRepo.Delete(uuidID)
}

func (uc *UserUseCase) GetUserMovementHistory(userID string) ([]*domain.BookMovementHistory, error) {
	uuidID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}
	return uc.movementRepo.GetByUserID(uuidID)
}
