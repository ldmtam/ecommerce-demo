package handlers

import (
	"github.com/Shopify/sarama"
	"github.com/ldmtam/ecommerce-demo/internal/models"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type repository interface {
	CreateProduct(name string, price uint) (*models.Product, error)
	GetProductByID(id uint) (*models.Product, error)
	GetProductByName(name string, limit uint) ([]*models.Product, error)
	GetCustomerActivities(id uint, limit uint) ([]*models.CustomerActivity, error)
	GetCustomerActivitiesByAction(id uint, action string, limit uint) ([]*models.CustomerActivity, error)
}

type handler struct {
	logger   *zap.Logger
	repo     repository
	producer sarama.SyncProducer
}

func New(logger *zap.Logger, repo repository) (*handler, error) {
	producer, err := initProducer(logger, viper.GetStringSlice("kafka.brokers"))
	if err != nil {
		return nil, err
	}

	return &handler{
		logger:   logger,
		repo:     repo,
		producer: producer,
	}, nil
}

func initProducer(logger *zap.Logger, brokers []string) (sarama.SyncProducer, error) {
	logger.Info("Creating kafka producer...")

	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)

	logger.Info("Successfully created kafka producer")

	return producer, err
}
