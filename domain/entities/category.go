package entities

type Category struct {
	ID    uint `gorm:"primaryKey,type:serial"`
	Image []byte
	Name  string `gorm:"not null"`
}
