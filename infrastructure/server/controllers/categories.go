package controllers

import (
	"kurs-server/domain/entities"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type createDto struct {
	Name string `json:"name" binding:"required"`
}

func (ctr *Controller) CreateCategory(c *gin.Context) {
	var dto createDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newCategory := &entities.Category{Name: dto.Name}
	err := ctr.cases.Categories().Create(newCategory)
	if err != nil {
		ctr.logger.Debug("Unable to create new category", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

type getDto struct {
	Name string `form:"name"`
}

func (ctr *Controller) GetCategories(c *gin.Context) {
	var dto getDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foundCategories := ctr.cases.Categories().GetMany(dto.Name)
	c.JSON(http.StatusOK, foundCategories)
}

type deleteDto struct {
	Name string `json:"name" binding:"required"`
}

func (ctr *Controller) DeleteCategory(c *gin.Context) {
	var dto deleteDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctr.cases.Categories().DeleteByName(dto.Name); err != nil {
		c.JSON(http.StatusBadRequest, "")
		return
	}

	c.Status(http.StatusOK)
}
