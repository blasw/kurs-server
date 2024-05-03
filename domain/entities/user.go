package entities

type User struct {
	ID           uint   `gorm:"primaryKey;type:serial"`
	Role         string `gorm:"default:User"`
	Username     string `gorm:"not null;unique"`
	Email        string `gorm:"not null"`
	Password     string `gorm:"not null"`
	RefreshToken string `gorm:"default:missing"`
}
