package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type getDetailsDto struct {
	ID uint `form:"id" binding:"required"`
}

func (ctr *Controller) GetDetails(c *gin.Context) {
	var dto getDetailsDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	details := ctr.cases.Details().GetForCategoryID(dto.ID)
	c.JSON(http.StatusOK, details)
}

type createValueDto struct {
	DetailID  uint   `json:"detail_id" binding:"required"`
	ProductID uint   `json:"product_id" binding:"required"`
	Value     string `json:"value" binding:"required"`
}

func (ctr *Controller) CreateValue(c *gin.Context) {
	var dto createValueDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctr.cases.Details().CreateValue(dto.DetailID, dto.ProductID, dto.Value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

type getDetailsValuesDto struct {
	CategoryID uint `form:"categoryid" binding:"required"`
}

func (ctr *Controller) GetDetailsValues(c *gin.Context) {
	var dto getDetailsValuesDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	detailsValues, err := ctr.cases.Details().GetUniqueValuesForCategoryID(dto.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, detailsValues)
}
