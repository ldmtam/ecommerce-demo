package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ldmtam/ecommerce-demo/internal/models"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

func (h *handler) GetCustomerActivitesByAction(c *gin.Context) {
	id := c.Param("id")
	customerID := cast.ToUint(id)
	if customerID == 0 {
		h.logger.Error("customer id is invalid", zap.String("id", id))
		c.JSON(http.StatusPreconditionFailed, gin.H{"error": "customer id is invalid"})
		return
	}

	action := c.Param("action_type")
	if action == "" || (action != models.CustomAction_SearchProduct && action != models.CustomAction_ViewProduct) {
		h.logger.Error("action is invalid", zap.String("id", id))
		c.JSON(http.StatusPreconditionFailed, gin.H{"error": "action is invalid"})
		return
	}

	activities, err := h.repo.GetCustomerActivitiesByAction(customerID, action, 20)
	if err != nil {
		h.logger.Error("Get customer activities by action failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Get customer activities by action failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": activities,
	})
}
