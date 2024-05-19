package controllers

import (
	"kurs-server/domain/entities"
	"kurs-server/structs"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type createReviewDto struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Text      string `json:"text" binding:"required"`
}

func (ctr *Controller) CreateReview(c *gin.Context) {
	var dto createReviewDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctr.logger.Debug("Creating review", zap.Any("dto", dto))

	userID := c.GetUint("ID")

	review := entities.Review{
		ProductID: dto.ProductID,
		UserID:    userID,
		Text:      dto.Text,
	}

	_, err := ctr.cases.Reviews().Create(&review)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

type getReviewsDto struct {
	ProductID uint `form:"product_id" binding:"required"`
}

type getReviewsResponse struct {
	Reviews []structs.FullReview `json:"reviews"`
}

func (ctr *Controller) GetReviews(c *gin.Context) {
	var dto getReviewsDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reviews := ctr.cases.Reviews().GetAllWithUsernames(dto.ProductID)

	if reviews == nil {
		reviews = []structs.FullReview{}
	}

	c.JSON(http.StatusOK, getReviewsResponse{Reviews: reviews})
}
