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

func (h *handler) SearchProductByName(c *gin.Context) {
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

	productName := c.Param("name")
	if productName == "" {
		h.logger.Error("product name is invalid", zap.String("name", productName))
		c.JSON(http.StatusPreconditionFailed, gin.H{"error": "product name is invalid"})
		return
	}

	products, err := h.repo.GetProductByName(productName, 20)
	if err != nil {
		h.logger.Error("Get products failed", zap.String("name", productName))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Get products failed"})
		return
	}

	// record view action asynchronously
	go h.publishSearchActivity(userID, products)

	c.JSON(http.StatusOK, gin.H{
		"data": products,
	})
}

func (h *handler) publishSearchActivity(userID uint, products []*models.Product) {
	productsBytes, _ := json.Marshal(products)
	activity := &models.CustomerActivity{
		UserID:    userID,
		CreatedAt: time.Now().UnixMilli(),
		Action:    models.CustomAction_SearchProduct,
		Data:      string(productsBytes),
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

	h.logger.Info("Recorded search product activity",
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
	)
}
