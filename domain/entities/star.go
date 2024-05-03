package entities

type Star struct {
	ID        uint    `gorm:"primaryKey"`
	Product   Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:CASCADE"`
	ProductID uint    `gorm:"not null"`
	User      User    `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	UserID    uint    `gorm:"not null"`
	Amount    int     `gorm:"not null"`
}
