package repos

import (
	"errors"
	"fmt"
	"kurs-server/domain/entities"
	"kurs-server/structs"
	"strings"

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

func (r *ProductRepo) Get(filters *structs.ProductFilters, limit int, offset int, sort string) ([]entities.Product, int64) {
	var products []entities.Product

	query := r.Storage.Model(&entities.Product{})
	var count int64

	if filters != nil {
		if filters.CategoryID != 0 {
			query = query.Joins("JOIN product_categories ON product_categories.product_id = products.id").Where("product_categories.category_id = ?", filters.CategoryID)
		}
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

	query.Count(&count)

	switch sort {
	case "id_desc":
		query.Limit(limit).Offset(offset).Order("id desc").Find(&products)

	case "id_asc":
		query.Limit(limit).Offset(offset).Order("id asc").Find(&products)

	case "price_desc":
		query.Limit(limit).Offset(offset).Order("price desc").Find(&products)

	case "price_asc":
		query.Limit(limit).Offset(offset).Order("price asc").Find(&products)
	}

	return products, count
}

func (r *ProductRepo) FilteredGet(offset int, limit int, sort string, search structs.ProductsSearch) ([]entities.Product, int64) {
	var products []entities.Product
	var count int64

	query := r.Storage.Model(&entities.Product{}).Preload("Categories").Preload("Details")
	query = r.applyFilters(query, search)

	query.Count(&count)

	switch sort {
	case "id_desc":
		query.Limit(limit).Offset(offset).Order("id desc").Find(&products)
	case "id_asc":
		query.Limit(limit).Offset(offset).Order("id asc").Find(&products)

	case "price_desc":
		query.Limit(limit).Offset(offset).Order("price desc").Find(&products)

	case "price_asc":
		query.Limit(limit).Offset(offset).Order("price asc").Find(&products)

	default:
		query.Limit(limit).Offset(offset).Order("id desc").Find(&products)
	}

	return products, count
}

func (r *ProductRepo) applyFilters(query *gorm.DB, search structs.ProductsSearch) *gorm.DB {
	if search.Name != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(search.Name)+"%")
	}

	if search.Brand != "" {
		query = query.Where("LOWER(brand) LIKE ?", "%"+strings.ToLower(search.Brand)+"%")
	}

	if search.MinPrice > 0 {
		query = query.Where("price >= ?", search.MinPrice)
	}

	if search.MaxPrice > 0 {
		query = query.Where("price <= ?", search.MaxPrice)
	}

	query = r.applyCategoryFilters(query, search.Categories)

	return query
}

func (r *ProductRepo) applyCategoryFilters(query *gorm.DB, categories []structs.CategoriesSearch) *gorm.DB {
	fmt.Println(categories)

	for _, cat := range categories {
		subQuery := r.Storage.Model(&entities.ProductCategory{}).Select("product_categories.product_id").Where("product_categories.category_id = ?", cat.ID)

		for i, detail := range cat.Details {
			alias := fmt.Sprintf("dv%d", i)
			joinCondition := fmt.Sprintf("JOIN detail_values %s ON product_categories.product_id = %s.product_id AND %s.detail_id = ? AND %s.value IN ?", alias, alias, alias, alias)
			subQuery = subQuery.Joins(joinCondition, detail.ID, detail.Values)
		}

		query = query.Where("products.id IN (?)", subQuery)
	}

	return query
}

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
	tx := r.Storage.Where("product_id = ?", productID).Delete(&entities.ProductCategory{})

	return tx.Error
}

func (r *ProductRepo) AssignCategoryByID(productID uint, categoryID uint) error {
	// assigning found category to a product
	newProductCategory := entities.ProductCategory{
		ProductID:  productID,
		CategoryID: categoryID,
	}
	return r.Storage.Create(&newProductCategory).Error
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
