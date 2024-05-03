package entities

type Detail struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}
