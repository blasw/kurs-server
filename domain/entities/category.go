package entities

type Category struct {
	ID   uint   `gorm:"primaryKey,type:serial"`
	Name string `gorm:"not null"`
}
