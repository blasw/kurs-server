package repos

import (
	"errors"
	"fmt"
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
func (r *DetailRepo) CreateValue(detailID uint, productID uint, value string) error {
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

func (r *DetailRepo) GetForCategoryID(categoryID uint) []entities.Detail {
	var details []entities.Detail
	err := r.Storage.Where("category_id = ?", categoryID).Find(&details).Error
	if err != nil {
		fmt.Println("unable to get details: ", err.Error())
		return nil
	}

	return details
}

func (r *DetailRepo) GetValue(detailID uint, productId uint) entities.DetailValue {
	var detailValue entities.DetailValue
	err := r.Storage.Where("detail_id = ? AND product_id = ?", detailID, productId).First(&detailValue).Error
	if err != nil {
		return entities.DetailValue{}
	}
	return detailValue
}

// gets all details for a given category
func (r *DetailRepo) GetForCategoryName(categoryName string) []entities.Detail {
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

func (r *DetailRepo) DeleteValues(productID uint) {
	r.Storage.Where("product_id = ?", productID).Delete(&entities.DetailValue{})
}

func (r *DetailRepo) EditValue(detailID uint, porductID uint, newValue string) {
	r.Storage.Model(&entities.DetailValue{}).Where("detail_id = ? AND product_id = ?", detailID, porductID).Update("value", newValue)
}
