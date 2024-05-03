package entities

type Order struct {
	ID           uint       `gorm:"primaryKey"`
	ProductID    uint       `gorm:"not null"`
	Product      Product    `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE"`
	Quantity     uint       `gorm:"not null;default:1"`
	OrderGroup   OrderGroup `gorm:"foreignKey:OrderGroupID;references:ID;constraint:OnDelete:CASCADE"`
	OrderGroupID uint       `gorm:"not null"`
}
