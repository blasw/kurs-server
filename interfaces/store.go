package interfaces

import "gorm.io/gorm"

type Store interface {
	DB() *gorm.DB
	Init(dsn string)
	AutomateRating()
}
