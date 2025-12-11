package users

import "gorm.io/gorm"

type User struct {
	gorm.Model        // GORM provides ID, CreatedAt, UpdatedAt, DeletedAt
	Username   string `gorm:"uniqueIndex;not null"`
	APIKey     string `gorm:"uniqueIndex;not null"`
}
