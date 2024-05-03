package entities

type Detail struct {
	Id   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}
