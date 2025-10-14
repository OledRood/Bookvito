package domain

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"` // Never expose password in JSON
	FullName  string    `json:"full_name"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Book represents a book available for exchange
type Book struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Author      string    `gorm:"not null" json:"author"`
	ISBN        string    `json:"isbn"`
	Description string    `gorm:"type:text" json:"description"`
	Condition   string    `json:"condition"` // new, good, fair, poor
	OwnerID     uint      `gorm:"not null" json:"owner_id"`
	Owner       User      `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Available   bool      `gorm:"default:true" json:"available"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Exchange represents a book exchange between users
type Exchange struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	RequesterID uint      `gorm:"not null" json:"requester_id"`
	Requester   User      `gorm:"foreignKey:RequesterID" json:"requester,omitempty"`
	OwnerID     uint      `gorm:"not null" json:"owner_id"`
	Owner       User      `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	BookID      uint      `gorm:"not null" json:"book_id"`
	Book        Book      `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Status      string    `gorm:"not null" json:"status"` // pending, accepted, rejected, completed
	Message     string    `gorm:"type:text" json:"message"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
