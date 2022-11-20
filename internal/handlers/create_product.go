package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CreateProductRequest struct {
	Name  string
	Price uint
}

func (h *handler) CreateProduct(c *gin.Context) {
	productInfo := &CreateProductRequest{}
	if err := c.ShouldBindJSON(productInfo); err != nil {
		h.logger.Error("Parsed product info failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.repo.CreateProduct(productInfo.Name, productInfo.Price)
	if err != nil {
		h.logger.Error("Create product failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Create product failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": product,
	})
}
