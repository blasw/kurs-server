package controllers

import (
	"kurs-server/domain/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateRating struct {
	UserID    uint `json:"user_id" binding:"required"`
	ProductID uint `json:"product_id" binding:"required"`
	Amount    int  `json:"amount" binding:"required"`
}

func (ctr *Controller) CreateRating(c *gin.Context) {
	var dto CreateRating
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	star := &entities.Star{
		UserID:    dto.UserID,
		ProductID: dto.ProductID,
		Amount:    dto.Amount,
	}

	rating := ctr.cases.Ratings().Get(dto.UserID, dto.ProductID)
	if rating != nil {
		ctr.cases.Ratings().Delete(rating.UserID, rating.ProductID)
	}

	ctr.cases.Ratings().Create(star)

	c.Status(http.StatusOK)
}

func (ctr *Controller) DeleteRating(c *gin.Context) {}
