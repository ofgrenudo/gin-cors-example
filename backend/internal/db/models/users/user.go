package users

import "gorm.io/gorm"

type User struct {
	gorm.Model        // GORM provides ID, CreatedAt, UpdatedAt, DeletedAt
	Email      string `gorm:"uniqueIndex;not null"`
	FirstName  string `gorm:"not null"`
	LastName   string `gorm:"not null"`
	Username   string `gorm:"uniqueIndex;not null"`
	APIKey     string `gorm:"uniqueIndex;not null"`
}
