package repos

import (
	"kurs-server/domain/entities"

	"gorm.io/gorm"
)

type RatingRepo struct {
	Storage *gorm.DB
}

func (r *RatingRepo) Create(new_product *entities.Star) {
	r.Storage.Create(new_product)
}

func (r *RatingRepo) Get(userID uint, productID uint) *entities.Star {
	var star entities.Star
	r.Storage.Where("user_id = ? AND product_id = ?", userID, productID).First(&star)
	return &star
}

func (r *RatingRepo) Delete(userID uint, productID uint) {
	r.Storage.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&entities.Star{})
}
