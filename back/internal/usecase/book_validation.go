package usecase

import (
	"bookvito/internal/domain"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	maxBookTitleLength       = 255
	maxBookAuthorLength      = 255
	maxBookDescriptionLength = 2000
	maxBookImageURLLength    = 500
)

func validateBookID(bookID uuid.UUID) error {
	if bookID == uuid.Nil {
		return domain.NewValidationError("Неверный идентификатор книги")
	}
	return nil
}

func validateTitle(title string) (string, error) {
	value := strings.TrimSpace(title)
	if value == "" {
		return "", domain.NewValidationError("Название книги не может быть пустым")
	}
	if utf8.RuneCountInString(value) > maxBookTitleLength {
		return "", domain.NewValidationError("Название книги не может быть длиннее 255 символов")
	}
	return value, nil
}

func validateAuthor(author string) (string, error) {
	value := strings.TrimSpace(author)
	if value == "" {
		return "", domain.NewValidationError("Автор книги не может быть пустым")
	}
	if utf8.RuneCountInString(value) > maxBookAuthorLength {
		return "", domain.NewValidationError("Автор книги не может быть длиннее 255 символов")
	}
	return value, nil
}

func validateDescription(description string) (string, error) {
	value := strings.TrimSpace(description)
	if utf8.RuneCountInString(value) > maxBookDescriptionLength {
		return "", domain.NewValidationError("Описание книги не может быть длиннее 2000 символов")
	}
	return value, nil
}

func validateCondition(condition domain.BookCondition) error {
	switch condition {
	case domain.ConditionExcellent, domain.ConditionGood, domain.ConditionBad:
		return nil
	default:
		return domain.NewValidationError("Неверное состояние книги")
	}
}

func (uc *BookUseCase) validateImageURL(imageURL string) (string, error) {
	value := strings.TrimSpace(imageURL)
	if value == "" {
		return "", nil
	}
	if utf8.RuneCountInString(value) > maxBookImageURLLength {
		return "", domain.NewValidationError("image_url не может быть длиннее 500 символов")
	}
	if uc.imageStorage == nil {
		return "", domain.NewInternalError("image storage is not configured", nil)
	}

	exists, err := uc.imageStorage.Exists(value)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", domain.NewValidationError("Файл изображения не найден на сервере")
	}

	return value, nil
}

func (uc *BookUseCase) validateLocation(locationID *uuid.UUID) error {
	if locationID == nil {
		return nil
	}
	if *locationID == uuid.Nil {
		return domain.NewValidationError("Неверный current_location_id")
	}
	if uc.locationRepo == nil {
		return nil
	}
	_, err := uc.locationRepo.GetByID(*locationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.NewNotFoundError("Локация не найдена")
		}
		return err
	}
	return nil
}

func (uc *BookUseCase) validateBookForCreate(book *domain.Book) error {
	if book == nil {
		return domain.NewValidationError("Книга не передана")
	}
	if book.OwnerID == uuid.Nil {
		return domain.NewValidationError("Неверный владелец книги")
	}

	title, err := validateTitle(book.Title)
	if err != nil {
		return err
	}
	author, err := validateAuthor(book.Author)
	if err != nil {
		return err
	}
	description, err := validateDescription(book.Description)
	if err != nil {
		return err
	}
	if err := validateCondition(book.Condition); err != nil {
		return err
	}
	imageURL, err := uc.validateImageURL(book.ImageURL)
	if err != nil {
		return err
	}
	if err := uc.validateLocation(book.CurrentLocationID); err != nil {
		return err
	}

	book.Title = title
	book.Author = author
	book.Description = description
	book.ImageURL = imageURL
	return nil
}
