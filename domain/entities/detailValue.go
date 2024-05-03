package entities

type DetailValue struct {
	ID        uint    `gorm:"primaryKey"`
	Value     string  `gorm:"not null"`
	ProductID uint    `gorm:"not null"`
	DetailID  uint    `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE"`
	Detail    Detail  `gorm:"foreignKey:DetailID;references:ID;constraint:OnDelete:CASCADE"`
}
