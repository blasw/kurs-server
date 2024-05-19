package repos

import (
	"kurs-server/domain/entities"
	"kurs-server/structs"

	"gorm.io/gorm"
)

type ReviewRepo struct {
	Storage *gorm.DB
}

func (r *ReviewRepo) Create(new_review *entities.Review) (uint, error) {
	tx := r.Storage.Create(new_review)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return new_review.ID, nil
}

func (r *ReviewRepo) GetAllWithUsernames(productID uint) []structs.FullReview {
	var reviews []entities.Review
	r.Storage.Preload("User").Where("product_id = ?", productID).Find(&reviews)

	var result []structs.FullReview

	for _, review := range reviews {
		result = append(result, structs.FullReview{
			ID:       review.ID,
			Username: review.User.Username,
			Text:     review.Text,
		})
	}

	return result
}
