package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type markOrderGroupDto struct {
	Mark string `form:"mark" binding:"required"`
	ID   uint   `form:"id" binding:"required"`
}

func (ctr *Controller) MarkOrderGroup(c *gin.Context) {
	var dto markOrderGroupDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctr.cases.Orders().MarkOrderGroup(dto.Mark, dto.ID)

	c.Status(http.StatusOK)
}

type deleteOrderGroupDto struct {
	ID uint `form:"id" binding:"required"`
}

func (ctr *Controller) DeleteOrderGroup(c *gin.Context) {
	var dto deleteOrderGroupDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctr.cases.Orders().DeleteOrderGroup(dto.ID)

	c.Status(http.StatusOK)
}
