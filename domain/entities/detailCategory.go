package entities

type DetailCategory struct {
	ID         uint     `gorm:"primaryKey;type:serial"`
	Detail     Detail   `gorm:"foreignKey:DetailID;references:ID;constraint:OnDelete:CASCADE"`
	Category   Category `gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:CASCADE"`
	DetailID   uint     `gorm:"not null"`
	CategoryID uint     `gorm:"not null"`
}
