package entities

type ProductCategory struct {
	ID         uint     `gorm:"primaryKey;type:serial"`
	Product    Product  `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE"`
	Category   Category `gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:CASCADE"`
	ProductID  uint     `gorm:"not null"`
	CategoryID uint     `gorm:"not null"`
}
