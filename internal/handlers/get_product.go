package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"github.com/ldmtam/ecommerce-demo/internal/models"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func (h *handler) GetProduct(c *gin.Context) {
	userIDCookie, err := c.Cookie("user_id")
	if err != nil {
		h.logger.Error("user authentication is invalid", zap.Error(err))
		c.JSON(http.StatusPreconditionFailed, gin.H{"error": "user authentication is invalid"})
		return
	}
	userID := cast.ToUint(userIDCookie)
	if userID == 0 {
		h.logger.Error("user authentication is invalid")
		c.JSON(http.StatusPreconditionFailed, gin.H{"error": "user authentication is invalid"})
		return
	}

	id := c.Param("id")
	productID := cast.ToUint(id)
	if productID == 0 {
		h.logger.Error("product id is invalid", zap.String("id", id))
		c.JSON(http.StatusPreconditionFailed, gin.H{"error": "product id is invalid"})
		return
	}

	product, err := h.repo.GetProductByID(productID)
	if err != nil {
		h.logger.Error("Get product failed", zap.Error(err), zap.String("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Get product failed"})
		return
	}

	// record view action asynchronously
	go h.publishViewActivity(userID, product)

	c.JSON(http.StatusOK, gin.H{
		"data": product,
	})
}

func (h *handler) publishViewActivity(userID uint, product *models.Product) {
	productBytes, _ := json.Marshal(product)
	activity := &models.CustomerActivity{
		UserID:    userID,
		CreatedAt: time.Now().UnixMilli(),
		Action:    models.CustomAction_ViewProduct,
		Data:      string(productBytes),
	}
	activityBytes, _ := json.Marshal(activity)

	partition, offset, err := h.producer.SendMessage(&sarama.ProducerMessage{
		Topic: viper.GetString("kafka.topic"),
		Value: sarama.StringEncoder(string(activityBytes)),
	})
	if err != nil {
		h.logger.Error("Produced message to kafka failed", zap.Error(err), zap.String("topic", viper.GetString("kafka.topic")))
		return
	}

	h.logger.Info("Recorded view product activity",
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
	)
}
