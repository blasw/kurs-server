package entities

import "time"

type OrderGroup struct {
	ID        uint
	Name      string
	Surname   string
	City      string
	Address   string
	Status    string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
	User      User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	UserID    uint
}
