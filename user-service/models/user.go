package models

type User struct {
	ID                  uint `gorm:"primaryKey"`
	OAuthProvider       string
	OAuthProviderUserID string `gorm:"not null;uniqueIndex"`
	Email               string
	Name                string
}
