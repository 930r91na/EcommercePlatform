package models

type User struct {
	ID                  uint   `gorm:"primaryKey"`
	OAuthProvider       string `gorm:"default:'local'"`
	OAuthProviderUserID string `gorm:"uniqueIndex"`
	Email               string
	Name                string
}
