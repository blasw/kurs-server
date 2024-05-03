package entities

import "time"

type OrderGroup struct {
	ID        uint      `gorm:"primaryKey"`
	Status    string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
	User      User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	UserID    uint
}
