package repos

import (
	"kurs-server/domain/entities"

	"gorm.io/gorm"
)

type CategoryRepo struct {
	Storage *gorm.DB
}

func (c *CategoryRepo) Create(newCategory *entities.Category) error {
	tx := c.Storage.Create(newCategory)
	return tx.Error
}

func (c *CategoryRepo) GetByName(name string) (*entities.Category, error) {
	var category entities.Category
	tx := c.Storage.Where("name = ?", name).First(&category)
	return &category, tx.Error
}

func (c *CategoryRepo) GetMany(name string) []entities.Category {
	var categories []entities.Category
	c.Storage.Where("LOWER(name) LIKE LOWER(?)", "%"+name+"%").Find(&categories)
	return categories
}

func (c *CategoryRepo) GetById(id int) (*entities.Category, error) {
	var category entities.Category
	tx := c.Storage.Where("id = ?", id).First(&category)
	return &category, tx.Error
}

// AssignCategory Todo: Might be incorrect
func (c *CategoryRepo) AssignCategory(productId uint, categoryId uint) error {
	tx := c.Storage.Create(&entities.ProductCategory{
		ProductID:  productId,
		CategoryID: categoryId,
	})

	return tx.Error
}

func (c *CategoryRepo) DeleteById(categoryId uint) error {
	return c.Storage.Delete(&entities.Category{ID: categoryId}).Error
}

func (c *CategoryRepo) DeleteByName(categoryName string) error {
	return c.Storage.Delete(&entities.Category{}).Where("name LIKE ?", categoryName).Error
}
