package usecase

import (
	"bookvito/internal/domain"
	"testing"
)

func TestBookUseCase_CreateBook_ValidatesImageViaStorage(t *testing.T) {
	bookRepo := newFakeBookRepo()
	movementRepo := &fakeMovementRepo{}
	exchangeRepo := &fakeExchangeRepo{}
	locationRepo := newFakeLocationRepo()
	imageStorage := newFakeImageStorage()
	imageURL := imageStorage.PublicURL("cover.jpg")
	imageStorage.saved[imageURL] = "image"

	location := &domain.Location{Name: "Central", Address: "Main street"}
	if err := locationRepo.Create(location); err != nil {
		t.Fatalf("seed location: %v", err)
	}

	uc := NewBookUseCase(bookRepo, movementRepo, exchangeRepo, locationRepo, &fakeMetadataProvider{}, imageStorage)
	book := &domain.Book{
		OwnerID:           mustUUID(),
		Title:             "Test Book",
		Author:            "Author",
		Description:       "Description",
		Condition:         domain.ConditionGood,
		ImageURL:          imageURL,
		CurrentLocationID: &location.ID,
	}

	if err := uc.CreateBook(book); err != nil {
		t.Fatalf("create book: %v", err)
	}
	if len(movementRepo.items) != 1 {
		t.Fatalf("expected movement entry, got %d", len(movementRepo.items))
	}
}

func TestBookUseCase_Request_OwnBookForbidden(t *testing.T) {
	bookRepo := newFakeBookRepo()
	movementRepo := &fakeMovementRepo{}
	exchangeRepo := &fakeExchangeRepo{}
	locationRepo := newFakeLocationRepo()
	imageStorage := newFakeImageStorage()
	ownerID := mustUUID()
	book := &domain.Book{
		ID:      mustUUID(),
		OwnerID: ownerID,
		Title:   "Owned Book",
		Author:  "Author",
		Status:  domain.BookAvailable,
	}
	if err := bookRepo.Create(book); err != nil {
		t.Fatalf("seed book: %v", err)
	}

	uc := NewBookUseCase(bookRepo, movementRepo, exchangeRepo, locationRepo, &fakeMetadataProvider{}, imageStorage)

	err := uc.Request(book.ID, ownerID, nil)
	if !hasErrorCode(err, domain.ErrorCodeForbidden) {
		t.Fatalf("expected forbidden error, got %v", err)
	}
}

func TestBookUseCase_SetBookImage_ReplacesOldImage(t *testing.T) {
	bookRepo := newFakeBookRepo()
	movementRepo := &fakeMovementRepo{}
	exchangeRepo := &fakeExchangeRepo{}
	locationRepo := newFakeLocationRepo()
	imageStorage := newFakeImageStorage()
	ownerID := mustUUID()
	oldURL := imageStorage.PublicURL("old.jpg")
	newURL := imageStorage.PublicURL("new.jpg")
	imageStorage.saved[oldURL] = "old"
	imageStorage.saved[newURL] = "new"

	book := &domain.Book{
		ID:       mustUUID(),
		OwnerID:  ownerID,
		Title:    "Book",
		Author:   "Author",
		Status:   domain.BookAvailable,
		ImageURL: oldURL,
	}
	if err := bookRepo.Create(book); err != nil {
		t.Fatalf("seed book: %v", err)
	}

	uc := NewBookUseCase(bookRepo, movementRepo, exchangeRepo, locationRepo, &fakeMetadataProvider{}, imageStorage)

	if err := uc.SetBookImage(book.ID, ownerID, false, newURL); err != nil {
		t.Fatalf("set image: %v", err)
	}

	updated, err := bookRepo.GetByID(book.ID)
	if err != nil {
		t.Fatalf("get book: %v", err)
	}
	if updated.ImageURL != newURL {
		t.Fatalf("expected image %q, got %q", newURL, updated.ImageURL)
	}
	if len(imageStorage.deleted) != 1 || imageStorage.deleted[0] != oldURL {
		t.Fatalf("expected old image delete, got %+v", imageStorage.deleted)
	}
}

func TestBookUseCase_DeleteBookImage_NoImageValidation(t *testing.T) {
	bookRepo := newFakeBookRepo()
	movementRepo := &fakeMovementRepo{}
	exchangeRepo := &fakeExchangeRepo{}
	locationRepo := newFakeLocationRepo()
	imageStorage := newFakeImageStorage()
	ownerID := mustUUID()

	book := &domain.Book{
		ID:      mustUUID(),
		OwnerID: ownerID,
		Title:   "Book",
		Author:  "Author",
		Status:  domain.BookAvailable,
	}
	if err := bookRepo.Create(book); err != nil {
		t.Fatalf("seed book: %v", err)
	}

	uc := NewBookUseCase(bookRepo, movementRepo, exchangeRepo, locationRepo, &fakeMetadataProvider{}, imageStorage)

	err := uc.DeleteBookImage(book.ID, ownerID, false)
	if !hasErrorCode(err, domain.ErrorCodeValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}
