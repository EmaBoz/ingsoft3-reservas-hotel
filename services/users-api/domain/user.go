package domain

import "time"

type User struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	Username       string     `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	Email          string     `gorm:"type:varchar(150);uniqueIndex;not null" json:"email"`
	Password       string     `gorm:"type:varchar(512);not null" json:"-"`
	FirstName      string     `gorm:"type:varchar(100)" json:"first_name"`
	LastName       string     `gorm:"type:varchar(100)" json:"last_name"`
	Role           UserRole   `gorm:"type:varchar(50)" json:"role"`
	Token          *string    `gorm:"type:varchar(512)" json:"token,omitempty"`
	TokenExpiresAt *time.Time `json:"token_expires_at,omitempty"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type UserRole string

const (
	RoleNormal UserRole = "normal"
	RoleAdmin  UserRole = "admin"
)
