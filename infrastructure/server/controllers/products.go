package controllers

import (
	"kurs-server/application/utility"
	"kurs-server/domain/entities"
	"kurs-server/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createNewProductDto struct {
	Brand       string   `json:"brand" binding:"required"`
	Name        string   `json:"name" binding:"required"`
	Price       float32  `json:"price" binding:"required"`
	Image       []byte   `json:"image"`
	Description string   `json:"description" binding:"required"`
	Categories  []string `json:"categories" binding:"required"`
}

func (ctr *Controller) CreateNewProduct(c *gin.Context) {
	var dto createNewProductDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(dto.Image) == 0 {
		defaultImage, _ := utility.GetDefaultImageBytes("./default_images/logo.png", "png")
		dto.Image = defaultImage
	}

	ctr.logger.Debug("Starting transcation")
	//starting transaction
	err := ctr.cases.Begin()
	if err != nil {
		ctr.logger.Debug("Failed to start transaction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctr.logger.Debug("Transaction started")

	//creating product
	newProduct := entities.Product{
		Name:        dto.Name,
		Brand:       dto.Brand,
		Price:       dto.Price,
		Image:       dto.Image,
		Description: dto.Description,
	}

	ctr.logger.Debug("Creating new product")

	productID, err := ctr.cases.Products().Create(&newProduct)
	if err != nil {
		ctr.logger.Debug("Failed to create new product")
		ctr.cases.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctr.logger.Debug("New product created, creating categories")

	for _, category := range dto.Categories {
		err = ctr.cases.Products().AssignCategoryByName(productID, category)
		if err != nil {
			ctr.logger.Debug("Failed to create categories for new product")
			ctr.cases.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctr.logger.Debug("Categories created, commiting")

	err = ctr.cases.Commit()
	if err != nil {
		ctr.logger.Debug("Failed to commit, rolling back")
		ctr.cases.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctr.logger.Debug("Success")

	c.JSON(http.StatusOK, gin.H{"productId": productID})
}

type deleteProductDto struct {
	ProductID uint `form:"id" binding:"required"`
}

func (ctr *Controller) DeleteProduct(c *gin.Context) {
	var dto deleteProductDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	// role := c.GetString("Role")
	// if role != "Admin" {
	// 	c.JSON(http.StatusForbidden, gin.H{"Message": "Permission denied"})
	// 	return
	// }

	err := ctr.cases.Products().DeleteById(dto.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message: ": "Success"})
}

// TODO: make filter by categories
type getProductsDto struct {
	Amount   int     `form:"amount" binding:"required"`
	Page     int     `form:"page" binding:"required"`
	Sort     string  `form:"sort"` // id_desc, id_asc, price_desc, price_asc
	Brand    string  `form:"brand"`
	Name     string  `form:"name"`
	MinPrice float32 `form:"minprice"`
	MaxPrice float32 `form:"maxprice"`
}

type fullProductResponse struct {
	ID          uint
	Name        string
	Brand       string
	Description string   `json:"desc"`
	Categories  []string `json:"categories"`
	Price       float32
	Rating      float32
	Image       []byte
}

func (ctr *Controller) GetProducts(c *gin.Context) {
	var dto getProductsDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if dto.Sort == "" {
		dto.Sort = "id_desc"
	}

	filters := structs.ProductFilters{
		Brand:    dto.Brand,
		Name:     dto.Name,
		MinPrice: dto.MinPrice,
		MaxPrice: dto.MaxPrice,
	}

	offset := dto.Amount * (dto.Page - 1)

	products, err := ctr.cases.Products().Get(&filters, dto.Amount, offset, dto.Sort)

	resp := []fullProductResponse{}

	for _, v := range products {
		categories := ctr.cases.Products().GetCategories(v.ID)
		categoriesResp := []string{}
		for _, category := range *categories {
			categoriesResp = append(categoriesResp, category.Name)
		}

		resp = append(resp, fullProductResponse{
			ID:          v.ID,
			Name:        v.Name,
			Brand:       v.Brand,
			Description: v.Description,
			Price:       v.Price,
			Rating:      v.Rating,
			Image:       v.Image,
			Categories:  categoriesResp,
		})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, resp)
}

type editProductDto struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Brand       string   `json:"brand"`
	Description string   `json:"desc"`
	Image       []byte   `json:"image"`
	Price       float32  `json:"price"`
	Categories  []string `json:"categories"`
}

func (ctr *Controller) EditProduct(c *gin.Context) {
	var dto editProductDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if c.GetString("Role") != "Admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"Message": "Access denied"})
		return
	}

	updatedCategories := []entities.Category{}
	for _, category := range dto.Categories {
		updatedCategories = append(updatedCategories, entities.Category{Name: category})
	}

	updatedProduct, err := ctr.cases.Products().Edit(dto.ID, dto.Name, dto.Brand, dto.Description, dto.Image, dto.Price, updatedCategories)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	temp := ctr.cases.Products().GetCategories(updatedProduct.ID)
	buf := []string{}
	for _, category := range *temp {
		buf = append(buf, category.Name)
	}

	resp := &fullProductResponse{
		ID:          updatedProduct.ID,
		Name:        updatedProduct.Name,
		Brand:       updatedProduct.Brand,
		Description: updatedProduct.Description,
		Image:       updatedProduct.Image,
		Price:       updatedProduct.Price,
		Categories:  buf,
	}

	c.JSON(http.StatusOK, resp)
}
