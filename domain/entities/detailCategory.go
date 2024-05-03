package entities

type DetailCategory struct {
	ID         uint     `gorm:"primaryKey"`
	DetailID   uint     `gorm:"not null"`
	CategoryID uint     `gorm:"not null"`
	Detail     Detail   `gorm:"foreignKey:DetailID;references:ID;constraint:OnDelete:CASCADE"`
	Category   Category `gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:CASCADE"`
}
