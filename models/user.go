package models

import "time"

// User model đại diện cho bảng users trong database
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`     // Không trả về trong JSON response
	Bio          *string   `json:"bio"`   // Có thể null
	Image        *string   `json:"image"` // Có thể null
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
