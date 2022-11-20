package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

func (h *handler) GetCustomerActivites(c *gin.Context) {
	id := c.Param("id")
	customerID := cast.ToUint(id)
	if customerID == 0 {
		h.logger.Error("customer id is invalid", zap.String("id", id))
		c.JSON(http.StatusPreconditionFailed, gin.H{"error": "customer id is invalid"})
		return
	}

	activities, err := h.repo.GetCustomerActivities(customerID, 20)
	if err != nil {
		h.logger.Error("Get customer activities failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Get customer activities failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": activities,
	})
}
