package entities

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey,type:serial"`
	Brand       string    `gorm:"not null"`
	Name        string    `gorm:"not null"`
	Image       []byte    `gorm:"not null"`
	Price       float32   `gorm:"not null"`
	Description string    `gorm:"not null" json:"desc"`
	Rating      float32   `gorm:"default:0.0"`
	CreatedAt   time.Time `gorm:"default:now()"`
}
