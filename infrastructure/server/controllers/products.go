package controllers

import (
	"kurs-server/application/utility"
	"kurs-server/domain/entities"
	"kurs-server/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createNewProductDto struct {
	Brand       string          `json:"brand" binding:"required"`
	Name        string          `json:"name" binding:"required"`
	Price       float32         `json:"price" binding:"required"`
	Image       []byte          `json:"image"`
	Description string          `json:"desc" binding:"required"`
	Categories  []categoriesDto `json:"categories"`
}

type categoriesDto struct {
	ID      uint         `json:"id"`
	Name    string       `json:"name"`
	Details []detailsDto `json:"details"`
}

type detailsDto struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (ctr *Controller) CreateNewProduct(c *gin.Context) {
	var dto createNewProductDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// using default image if no image provided
	if len(dto.Image) == 0 {
		defaultImage, _ := utility.GetDefaultImageBytes("./default_images/logo.png", "png")
		dto.Image = defaultImage
	}

	//starting transaction
	err := ctr.cases.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//creating product
	newProduct := entities.Product{
		Name:        dto.Name,
		Brand:       dto.Brand,
		Price:       dto.Price,
		Image:       dto.Image,
		Description: dto.Description,
	}

	//creating product

	productID, err := ctr.cases.Products().Create(&newProduct)
	if err != nil {
		ctr.logger.Debug("Failed to create new product")
		_ = ctr.cases.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//binding categories
	for _, category := range dto.Categories {
		err = ctr.cases.Products().AssignCategoryByID(productID, category.ID)
		if err != nil {
			ctr.logger.Debug("Failed to create categories for new product")
			_ = ctr.cases.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		//creating details values
		for _, detail := range category.Details {
			_ = ctr.cases.Details().CreateValue(detail.ID, productID, detail.Value)
		}
	}

	err = ctr.cases.Commit()
	if err != nil {
		ctr.logger.Debug("Failed to commit, rolling back")
		_ = ctr.cases.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, productID)
}

type deleteProductDto struct {
	ProductID uint `form:"id" binding:"required"`
}

func (ctr *Controller) DeleteProduct(c *gin.Context) {
	var dto deleteProductDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	err := ctr.cases.Products().DeleteById(dto.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message: ": "Success"})
}

func (ctr *Controller) FilteredGetProducts(c *gin.Context) {
	var dto structs.ProductsSearch
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetUint("ID")

	limit := dto.Amount
	offset := (dto.Page - 1) * dto.Amount

	products, count := ctr.cases.Products().FilteredGet(offset, limit, dto.Sort, dto)

	prods := []getProductInfoResp{}

	for _, product := range products {
		prods = append(prods, ctr.getFullProduct(product.ID, userID))
	}

	resp := amountResp{
		Total:    int(count),
		Products: prods,
	}

	c.JSON(http.StatusOK, resp)
}

type getProductsDto struct {
	Amount     int     `form:"amount" binding:"required"`
	Page       int     `form:"page" binding:"required"`
	Sort       string  `form:"sort"` // id_desc, id_asc, price_desc, price_asc
	Brand      string  `form:"brand"`
	Name       string  `form:"name"`
	CategoryID uint    `form:"categoryid"`
	MinPrice   float32 `form:"minprice"`
	MaxPrice   float32 `form:"maxprice"`
}

type amountResp struct {
	Total    int                  `json:"total"`
	Products []getProductInfoResp `json:"products"`
}

type getProductInfoResp struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Brand       string         `json:"brand"`
	Description string         `json:"desc"`
	Image       []byte         `json:"image"`
	Price       float32        `json:"price"`
	Rating      float32        `json:"rating"`
	UserRating  float32        `json:"user_rating"`
	Categories  []categoryResp `json:"categories"`
}

type categoryResp struct {
	ID      uint          `json:"id"`
	Name    string        `json:"name"`
	Details []detailsResp `json:"details"`
}

type detailsResp struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
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
		Brand:      dto.Brand,
		Name:       dto.Name,
		CategoryID: dto.CategoryID,
		MinPrice:   dto.MinPrice,
		MaxPrice:   dto.MaxPrice,
	}

	offset := dto.Amount * (dto.Page - 1)

	products, count := ctr.cases.Products().Get(&filters, dto.Amount, offset, dto.Sort)

	prods := []getProductInfoResp{}

	userID, _ := c.Get("ID")

	for _, v := range products {
		fullProduct := ctr.getFullProduct(v.ID, userID.(uint))
		prods = append(prods, fullProduct)
	}

	resp := amountResp{
		Total:    int(count),
		Products: prods,
	}

	c.JSON(http.StatusOK, resp)
}

type editProductDto struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Brand       string          `json:"brand"`
	Description string          `json:"desc"`
	Image       []byte          `json:"image"`
	Price       float32         `json:"price"`
	Categories  []categoriesDto `json:"categories"`
}

func (ctr *Controller) EditProduct(c *gin.Context) {
	var dto editProductDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if len(dto.Categories) > 0 {
		ctr.cases.Details().DeleteValues(dto.ID)
	}

	updatedCategories := []entities.Category{}
	for _, category := range dto.Categories {
		updatedCategories = append(updatedCategories, entities.Category{ID: category.ID, Name: category.Name})

		for _, detail := range category.Details {
			ctr.cases.Details().CreateValue(detail.ID, dto.ID, detail.Value)
		}
	}

	_, err := ctr.cases.Products().Edit(dto.ID, dto.Name, dto.Brand, dto.Description, dto.Image, dto.Price, updatedCategories)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	resp := ctr.getFullProduct(dto.ID, 0)

	c.JSON(http.StatusOK, resp)
}

type getProductInfoDto struct {
	ID uint `form:"id" binding:"required"`
}

// TODO: should return reviews as well
func (ctr *Controller) GetProductInfo(c *gin.Context) {
	var dto getProductInfoDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetUint("ID")

	response := ctr.getFullProduct(dto.ID, userID)

	c.JSON(http.StatusOK, response)
}

func (ctr *Controller) getFullProduct(ID uint, userID uint) getProductInfoResp {
	// getting product
	product := ctr.cases.Products().GetProductByID(ID)

	// getting categories
	productCategories := ctr.cases.Products().GetCategories(ID)

	categories := []categoryResp{}

	// getting details
	for _, category := range *productCategories {
		detailsForCategory := ctr.cases.Details().GetForCategoryID(category.ID)
		tempCategory := categoryResp{
			ID:      category.ID,
			Name:    category.Name,
			Details: []detailsResp{},
		}
		for _, detail := range detailsForCategory {
			value := ctr.cases.Details().GetValue(detail.ID, ID).Value
			tempCategory.Details = append(tempCategory.Details, detailsResp{ID: detail.ID, Name: detail.Name, Value: value})
		}
		categories = append(categories, tempCategory)
	}

	userRating := ctr.cases.Ratings().Get(userID, ID)
	if userRating == nil {
		userRating = &entities.Star{
			UserID:    userID,
			ProductID: ID,
			Amount:    0,
		}
	}

	response := getProductInfoResp{
		ID:          product.ID,
		Name:        product.Name,
		Brand:       product.Brand,
		Description: product.Description,
		Image:       product.Image,
		Price:       product.Price,
		Rating:      product.Rating,
		UserRating:  float32(userRating.Amount),
		Categories:  categories,
	}

	return response
}
