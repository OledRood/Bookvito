package usecase

import (
	"bookvito/internal/domain"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func TestUserUseCase_RegisterUser_DuplicateEmailConflict(t *testing.T) {
	userRepo := newFakeUserRepo()
	movementRepo := &fakeMovementRepo{}
	existing := &domain.User{Email: "dup@example.com", Password: "hash", Name: "Dup"}
	if err := userRepo.Create(existing); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	uc := NewUserUseCase(userRepo, movementRepo, "secret")

	_, err := uc.RegisterUser("dup@example.com", "password123", "Another")
	if !hasErrorCode(err, domain.ErrorCodeConflict) {
		t.Fatalf("expected conflict error, got %v", err)
	}
}

func TestUserUseCase_LoginUser_SuccessGeneratesTokens(t *testing.T) {
	userRepo := newFakeUserRepo()
	movementRepo := &fakeMovementRepo{}
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{Email: "reader@example.com", Password: string(passwordHash), Name: "Reader", Role: domain.RoleUser}
	if err := userRepo.Create(user); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	uc := NewUserUseCase(userRepo, movementRepo, "secret")

	tokens, err := uc.LoginUser("reader@example.com", "password123")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		t.Fatalf("expected token pair, got %+v", tokens)
	}
}

func TestUserUseCase_RefreshToken_ExpiredReturnsError(t *testing.T) {
	userRepo := newFakeUserRepo()
	movementRepo := &fakeMovementRepo{}
	user := &domain.User{
		Email:                 "reader@example.com",
		Name:                  "Reader",
		RefreshToken:          "refresh-token",
		RefreshTokenExpiresAt: time.Now().Add(-time.Hour),
	}
	if err := userRepo.Create(user); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	uc := NewUserUseCase(userRepo, movementRepo, "secret")

	_, err := uc.RefreshToken("refresh-token")
	if !errContains(err, "истёк") {
		t.Fatalf("expected expired token error, got %v", err)
	}
}

func TestUserUseCase_Logout_ClearsRefreshToken(t *testing.T) {
	userRepo := newFakeUserRepo()
	movementRepo := &fakeMovementRepo{}
	user := &domain.User{
		Email:                 "reader@example.com",
		Name:                  "Reader",
		RefreshToken:          "refresh-token",
		RefreshTokenExpiresAt: time.Now().Add(time.Hour),
	}
	if err := userRepo.Create(user); err != nil {
		t.Fatalf("seed user: %v", err)
	}

	uc := NewUserUseCase(userRepo, movementRepo, "secret")

	if err := uc.Logout(user.ID.String()); err != nil {
		t.Fatalf("logout: %v", err)
	}

	updated, err := userRepo.GetByID(user.ID)
	if err != nil {
		t.Fatalf("read user: %v", err)
	}
	if updated.RefreshToken != "" {
		t.Fatalf("expected refresh token to be cleared, got %q", updated.RefreshToken)
	}
}
