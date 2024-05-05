package entities

type Detail struct {
	ID         uint     `gorm:"primaryKey"`
	Name       string   `gorm:"not null"`
	CategoryID uint     `gorm:"not null"`
	Category   Category `gorm:"foreignKey:CategoryID;references:ID;constraint:OnDelete:CASCADE"`
}
