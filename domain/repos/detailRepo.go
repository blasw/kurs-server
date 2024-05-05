package repos

import (
	"errors"
	"kurs-server/domain/entities"

	"gorm.io/gorm"
)

type DetailRepo struct {
	Storage *gorm.DB
}

// creates a plain detail with binded category
func (r *DetailRepo) Create(detailName string, categoryName string) error {
	var searchResult entities.Category
	err := r.Storage.Where("name = ?", categoryName).First(&searchResult).Error
	if err != nil {
		return errors.New("category not found")
	}

	new_detail := entities.Detail{Name: detailName, CategoryID: searchResult.ID}
	err = r.Storage.Create(&new_detail).Error
	if err != nil {
		return err
	}

	return nil
}

// assigns a value to a detailValue
func (r *DetailRepo) AssignValue(detailID uint, productID uint, value string) error {
	detailValue := entities.DetailValue{
		Value:     value,
		ProductID: productID,
		DetailID:  detailID,
	}
	if err := r.Storage.Create(&detailValue).Error; err != nil {
		return err
	}

	return nil
}

// gets all details for a given category
func (r *DetailRepo) GetForCategory(categoryName string) []entities.Detail {
	var searchCategory entities.Category
	err := r.Storage.Where("name = ?", categoryName).First(&searchCategory).Error
	if err != nil {
		return nil
	}

	var details []entities.Detail

	err = r.Storage.Where("category_id = ?", searchCategory.ID).Find(&details).Error
	if err != nil {
		return nil
	}

	return details
}
