package controllers

import (
	"kurs-server/domain/entities"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type createOrderDto struct {
	ProductID uint `form:"product_id" binding:"required"`
}

// Guarded
func (ctr *Controller) CreateOrder(c *gin.Context) {
	var dto createOrderDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("ID")

	orderGroup := ctr.cases.Orders().GetPendingOrderGroup(userID)

	order := entities.Order{
		OrderGroupID: orderGroup.ID,
		ProductID:    dto.ProductID,
	}

	ctr.cases.Orders().CreateOrder(&order)

	c.Status(http.StatusOK)
}

func (ctr *Controller) DeleteOrder(c *gin.Context) {
	var dto createOrderDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("ID")

	orderGroup := ctr.cases.Orders().GetPendingOrderGroup(userID)

	ctr.cases.Orders().DeleteOrder(orderGroup.ID, dto.ProductID)

	c.Status(http.StatusOK)
}

type cartProductResp struct {
	ID          uint    `json:"id"`
	Brand       string  `json:"brand"`
	Name        string  `json:"name"`
	Image       []byte  `json:"image"`
	Price       float32 `json:"price"`
	Description string  `json:"desc"`
	Quantity    uint    `json:"quantity"`
}

func (ctr *Controller) GetCart(c *gin.Context) {
	userID := c.GetUint("ID")

	orderGroup := ctr.cases.Orders().GetPendingOrderGroup(userID)

	orders := ctr.cases.Orders().GetOrdersByOrderGroupID(orderGroup.ID)

	products := []cartProductResp{}

	var product *entities.Product

	for _, order := range orders {
		product = ctr.cases.Products().GetProductByID(order.ProductID)
		products = append(products, cartProductResp{
			ID:          product.ID,
			Brand:       product.Brand,
			Name:        product.Name,
			Image:       product.Image,
			Price:       product.Price,
			Description: product.Description,
			Quantity:    order.Quantity,
		})
	}

	c.JSON(http.StatusOK, products)
}

type confirmOrderGroupDto struct {
	Name    string `json:"name" binding:"required"`
	Surname string `json:"surname" binding:"required"`
	City    string `json:"city" binding:"required"`
	Address string `json:"address" binding:"required"`
}

func (ctr *Controller) ConfirmOrderGroup(c *gin.Context) {
	var dto confirmOrderGroupDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("ID")

	orderGroup := ctr.cases.Orders().GetPendingOrderGroup(userID)

	ctr.cases.Orders().ProcessOrderGroup(orderGroup.ID, dto.Name, dto.Surname, dto.City, dto.Address)

	c.Status(http.StatusOK)
}

type getOrdersGroupsResp struct {
	ID        uint   `json:"id"`
	OrderDate string `json:"order_date"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	City      string `json:"city"`
	Address   string `json:"address"`
	Status    string `json:"status"`
}

func (ctr *Controller) GetOrderGroups(c *gin.Context) {
	userID := c.GetUint("ID")

	orderGroups := ctr.cases.Orders().GetOrderGroupsByUserID(userID)

	orderGroupsResp := []getOrdersGroupsResp{}

	for _, orderGroup := range orderGroups {
		if orderGroup.Status == "pending" {
			continue
		}

		orderGroupsResp = append(orderGroupsResp, getOrdersGroupsResp{
			ID:        orderGroup.ID,
			OrderDate: orderGroup.CreatedAt.String(),
			Name:      orderGroup.Name,
			Surname:   orderGroup.Surname,
			City:      orderGroup.City,
			Address:   orderGroup.Address,
			Status:    orderGroup.Status,
		})
	}

	c.JSON(http.StatusOK, orderGroupsResp)
}

type getAllOrderGroupsDto struct {
	Amount int `form:"amount"`
	Page   int `form:"page"`
}

func (ctr *Controller) GetAllOrderGroups(c *gin.Context) {
	var dto getAllOrderGroupsDto

	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit := dto.Amount
	offset := (dto.Page - 1) * dto.Amount

	orderGroups, count := ctr.cases.Orders().GetAllOrderGroups(offset, limit)

	ctr.logger.Debug("aboba", zap.Int("count", len(orderGroups)))

	orderGroupsResp := []getOrdersGroupsResp{}

	for _, orderGroup := range orderGroups {
		if orderGroup.Status == "pending" {
			continue
		}

		orderGroupsResp = append(orderGroupsResp, getOrdersGroupsResp{
			ID:        orderGroup.ID,
			OrderDate: orderGroup.CreatedAt.String(),
			Name:      orderGroup.Name,
			Surname:   orderGroup.Surname,
			City:      orderGroup.City,
			Address:   orderGroup.Address,
			Status:    orderGroup.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"order_groups": orderGroupsResp,
		"amount":       count,
	})
}
