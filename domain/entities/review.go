package entities

import "time"

type Review struct {
	ID        uint      `gorm:"primaryKey"`
	Text      string    `gorm:"not null"`
	Product   Product   `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE"`
	ProductID uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	UserID    uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:now()"`
}
