package repos

import (
	"errors"
	"fmt"
	"kurs-server/domain/entities"
	"kurs-server/structs"

	"gorm.io/gorm"
)

type ProductRepo struct {
	Storage *gorm.DB
}

func (r *ProductRepo) Create(new_product *entities.Product) (uint, error) {
	tx := r.Storage.Create(new_product)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return new_product.ID, nil
}

func (r *ProductRepo) DeleteById(productID uint) error {
	var searchResult []entities.Product
	r.Storage.Where("id = ?", productID).Find(&searchResult)

	if len(searchResult) != 1 {
		return errors.New("product not found")
	}

	tx := r.Storage.Delete(&searchResult[0])
	return tx.Error
}

func (r *ProductRepo) Get(filters *structs.ProductFilters, limit int, offset int, sort string) ([]entities.Product, error) {
	var products []entities.Product

	query := r.Storage.Model(&entities.Product{})

	if filters != nil {
		if filters.Brand != "" {
			query = query.Where("LOWER(brand) LIKE LOWER(?)", "%"+filters.Brand+"%")
		}
		if filters.Name != "" {
			query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+filters.Name+"%")
		}
		if filters.MinPrice > 0 {
			query = query.Where("price >= ?", filters.MinPrice)
		}
		if filters.MaxPrice > 0 {
			query = query.Where("price <= ?", filters.MaxPrice)
		}
	}

	var err error

	switch sort {
	case "id_desc":
		err = query.Limit(limit).Offset(offset).Order("id desc").Find(&products).Error

	case "id_asc":
		err = query.Limit(limit).Offset(offset).Order("id asc").Find(&products).Error

	case "price_desc":
		err = query.Limit(limit).Offset(offset).Order("price desc").Find(&products).Error

	case "price_asc":
		err = query.Limit(limit).Offset(offset).Order("price asc").Find(&products).Error
	}

	if err != nil {
		return nil, err
	}

	return products, nil
}

// Edit TODO: Might work wrong
func (r *ProductRepo) Edit(ID uint, Name string, Brand string, Description string, Image []byte, Price float32, Categories []entities.Category) (*entities.Product, error) {
	var product entities.Product

	if err := r.Storage.First(&product, ID).Error; err != nil {
		return nil, err
	}

	if Name != "" {
		product.Name = Name
	}
	if Brand != "" {
		product.Brand = Brand
	}
	if Description != "" {
		product.Description = Description
	}
	if len(Image) > 0 {
		product.Image = Image
	}
	if Price != 0 {
		product.Price = Price
	}

	if err := r.Storage.Save(&product).Error; err != nil {
		return nil, err
	}

	if len(Categories) > 0 {
		err := r.DeleteCategories(ID)
		if err != nil {
			return nil, err
		}

		for _, category := range Categories {
			r.AssignCategoryByName(ID, category.Name)
		}
	}

	return &product, nil
}

func (r *ProductRepo) DeleteCategories(productID uint) error {
	tx := r.Storage.Delete(&entities.ProductCategory{ProductID: productID})

	return tx.Error
}

func (r *ProductRepo) AssignCategoryByName(productID uint, categoryName string) error {
	//finding category by name
	var category entities.Category
	if err := r.Storage.First(&category, entities.Category{Name: categoryName}).Error; err != nil {
		return err
	}

	// assigning found category to a product
	newProductCategory := entities.ProductCategory{
		ProductID:  productID,
		CategoryID: category.ID,
	}

	if err := r.Storage.Create(&newProductCategory).Error; err != nil {
		return err
	}

	return nil
}

func (r *ProductRepo) GetCategories(productID uint) *[]entities.Category {
	var categoriesEntries []entities.ProductCategory
	tx := r.Storage.Where("product_id = ?", productID).Find(&categoriesEntries)

	if tx.Error != nil && categoriesEntries == nil {
		return nil
	}

	fmt.Println(categoriesEntries)

	categories := []entities.Category{}
	for _, entry := range categoriesEntries {
		var category entities.Category
		r.Storage.Where("id = ?", entry.CategoryID).First(&category)
		categories = append(categories, category)
	}

	return &categories
}

func (r *ProductRepo) GetProductByID(ID uint) *entities.Product {
	var product entities.Product
	r.Storage.Where("id = ?", ID).First(&product)
	return &product
}
