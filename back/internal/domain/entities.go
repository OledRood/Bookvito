package domain

import (
	"time"

	"github.com/google/uuid"
)

// Enum types
type BookStatus string

const (
	BookAvailable BookStatus = "available"
	BookRequested BookStatus = "requested"
	BookBorrowed  BookStatus = "borrowed"
	BookArchived  BookStatus = "archived"
	BookDeleted   BookStatus = "deleted"
)

type BookCondition string

const (
	ConditionExcellent BookCondition = "excellent"
	ConditionGood      BookCondition = "good"
	ConditionBad       BookCondition = "bad"
)

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleModer UserRole = "moder"
	RoleAdmin UserRole = "admin"
)

type ExchangeStatus string

const (
	ExchangeRequested ExchangeStatus = "requested"
	ExchangeBorrowed  ExchangeStatus = "borrowed"
	ExchangeReturned  ExchangeStatus = "returned"
	ExchangeCancelled ExchangeStatus = "cancelled"
	ExchangeOverdue   ExchangeStatus = "overdue"
)

// Location represents a pickup/return location (Пункт выдачи)
type Location struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name    string    `gorm:"not null" json:"name"`    // Название (например Библиотека №4)
	Address string    `gorm:"not null" json:"address"` // Адрес (ОБЯЗАТЕЛЬНО)
	Books   []Book    `gorm:"foreignKey:CurrentLocationID" json:"books,omitempty"`
}

// User represents a user in the system (Пользователь)
type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email     string     `gorm:"unique;not null" json:"email"`                // Почта
	Password  string     `gorm:"not null" json:"-"`                           // Пароль (хранится в виде хэша)
	Name      string     `gorm:"not null" json:"name"`                        // Имя
	Role      UserRole   `gorm:"type:varchar(20);default:'user'" json:"role"` // Роль
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`            // Дата создания
	Exchanges []Exchange `gorm:"foreignKey:UserID" json:"exchanges,omitempty"`
	Reviews   []Review   `gorm:"foreignKey:UserID" json:"reviews,omitempty"`

	RefreshToken          string    `json:"-"` // Поле для Refresh токена
	RefreshTokenExpiresAt time.Time `json:"-"` // Время жизни Refresh токена
}

// Book represents a book (Книга)
type Book struct {
	ID                uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OwnerID           uuid.UUID     `gorm:"type:uuid;not null" json:"owner_id"` // ID владельца книги
	Owner             *User         `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Title             string        `gorm:"not null" json:"title"`                              // Название
	Author            string        `gorm:"not null" json:"author"`                             // Автор
	Description       string        `gorm:"type:text" json:"description"`                       // Описание
	Condition         BookCondition `gorm:"type:varchar(20);default:'good'" json:"condition"`   // Состояние
	Status            BookStatus    `gorm:"type:varchar(20);default:'available'" json:"status"` // Бронь (состояние)
	CurrentLocationID *uuid.UUID    `gorm:"type:uuid" json:"current_location_id"`               // Текущая позиция (ссылка на пункт выдачи)
	CurrentLocation   *Location     `gorm:"foreignKey:CurrentLocationID" json:"current_location,omitempty"`
	ImageURL          string        `json:"image_url"`                        // Картинка
	UpdatedAt         time.Time     `gorm:"autoUpdateTime" json:"updated_at"` // Последнее изменение
	CreatedAt         time.Time     `gorm:"autoCreateTime" json:"created_at"`
	Reviews           []Review      `gorm:"foreignKey:BookID" json:"reviews,omitempty"` // Отзывы
	Exchanges         []Exchange    `gorm:"foreignKey:BookID" json:"exchanges,omitempty"`
}

type BookSummary struct {
	ID       uuid.UUID `json:"id" db:"id"`
	ImageURL string    `json:"image_url" db:"image_url"`
	Title    string    `json:"title" db:"title"`
	Author   string    `json:"author" db:"author"`
}

// Exchange represents a book reservation/borrowing (Бронирование)
type Exchange struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"` // Кто забронировал
	User       User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BookID     uuid.UUID      `gorm:"type:uuid;not null" json:"book_id"` // Книга
	Book       Book           `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Status     ExchangeStatus `gorm:"type:varchar(20);default:'requested'" json:"status"` // Статус
	BookedAt   time.Time      `gorm:"autoCreateTime" json:"booked_at"`                    // Когда забронировано
	ExpiresAt  *time.Time     `json:"expires_at"`                                         // Время, до которого бронь действительна
	LocationID *uuid.UUID     `gorm:"type:uuid" json:"location_id"`                       // Пункт выдачи
	Location   *Location      `gorm:"foreignKey:LocationID" json:"location,omitempty"`
}

// Review represents a book review (Отзывы)
type Review struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BookID    uuid.UUID `gorm:"type:uuid;not null" json:"book_id"` // Книга
	Book      Book      `gorm:"foreignKey:BookID" json:"book,omitempty"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"` // Пользователь
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Rating    int16     `gorm:"type:smallint;not null" json:"rating"` // Оценка (маленький int)
	Text      string    `gorm:"type:text" json:"text"`                // Текст отзыва
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`     // Дата создания
}

// BookMovementHistory represents the history of book movements (История перемещений книги)
type BookMovementHistory struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BookID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"book_id"` // Книга
	Book           *Book      `gorm:"foreignKey:BookID" json:"book,omitempty"`
	FromLocationID *uuid.UUID `gorm:"type:uuid" json:"from_location_id"` // Откуда (может быть NULL при первой регистрации)
	FromLocation   *Location  `gorm:"foreignKey:FromLocationID" json:"from_location,omitempty"`
	ToLocationID   *uuid.UUID `gorm:"type:uuid" json:"to_location_id"` // Куда (может быть NULL если книга взята пользователем)
	ToLocation     *Location  `gorm:"foreignKey:ToLocationID" json:"to_location,omitempty"`
	ExchangeID     *uuid.UUID `gorm:"type:uuid;index" json:"exchange_id"` // Связанное бронирование (если есть)
	Exchange       *Exchange  `gorm:"foreignKey:ExchangeID" json:"exchange,omitempty"`
	UserID         *uuid.UUID `gorm:"type:uuid" json:"user_id"` // Пользователь, который инициировал перемещение
	User           *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Action         string     `gorm:"type:varchar(50);not null" json:"action"` // Действие: created, moved, borrowed, returned, etc.
	Notes          string     `gorm:"type:text" json:"notes"`                  // Дополнительные заметки
	PreviousStatus BookStatus `gorm:"type:varchar(20)" json:"previous_status"` // Предыдущий статус книги
	NewStatus      BookStatus `gorm:"type:varchar(20)" json:"new_status"`      // Новый статус книги
	CreatedAt      time.Time  `gorm:"autoCreateTime;index" json:"created_at"`  // Время перемещения
}
