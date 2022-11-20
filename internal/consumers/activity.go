package consumers

import (
	"context"
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/ldmtam/ecommerce-demo/internal/models"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type repository interface {
	CreateCustomerActivity(userID uint, createdAt int64, action, data string) (*models.CustomerActivity, error)
}

type ActivityConsumer struct {
	logger   *zap.Logger
	repo     repository
	client   sarama.ConsumerGroup
	ready    chan bool
	ctx      context.Context
	cancelFn context.CancelFunc
}

func NewActivityConsumer(logger *zap.Logger, repo repository) (*ActivityConsumer, error) {
	client, err := initConsumer(logger, viper.GetStringSlice("kafka.brokers"))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &ActivityConsumer{
		logger:   logger,
		repo:     repo,
		client:   client,
		ready:    make(chan bool),
		ctx:      ctx,
		cancelFn: cancel,
	}, nil
}

func (c *ActivityConsumer) Start() error {
	go func() {
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := c.client.Consume(
				c.ctx,
				[]string{viper.GetString("kafka.topic")},
				&consumerHandler{c: c}); err != nil {
				c.logger.Error("Error from consumer", zap.Error(err))
			}
			// check if context was cancelled, signaling that the consumer should stop
			if c.ctx.Err() != nil {
				return
			}
			c.ready = make(chan bool)
		}
	}()

	<-c.ready

	return nil
}

type consumerHandler struct {
	c *ActivityConsumer
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (handler *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(handler.c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (handler *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (handler *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message := <-claim.Messages():
			handler.c.logger.Info("Message claimed",
				zap.Time("timestamp", message.Timestamp),
				zap.String("topic", message.Topic),
				zap.Int32("partition", message.Partition),
				zap.Int64("offset", message.Offset))

			customerActivity := &models.CustomerActivity{}
			if err := json.Unmarshal(message.Value, customerActivity); err != nil {
				handler.c.logger.Error("Parse json failed",
					zap.Error(err),
					zap.String("topic", message.Topic),
					zap.Int32("partition", message.Partition),
					zap.Int64("offset", message.Offset),
					zap.String("message", string(message.Value)))
				continue
			}

			if _, err := handler.c.repo.CreateCustomerActivity(
				customerActivity.UserID,
				customerActivity.CreatedAt,
				customerActivity.Action,
				customerActivity.Data,
			); err != nil {
				handler.c.logger.Error("Create customer activity failed",
					zap.Error(err),
					zap.Reflect("customer activity", customerActivity))
				continue
			}

			session.MarkMessage(message, "")

		// Should return when `session.Context()` is done.
		// If not, will r`aise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

func (c *ActivityConsumer) Stop() {
	c.cancelFn()
	<-c.ctx.Done()
}

func initConsumer(logger *zap.Logger, brokers []string) (sarama.ConsumerGroup, error) {
	logger.Info("Creating kafka consumer...")

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V1_0_0_0
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(
		viper.GetStringSlice("kafka.brokers"),
		viper.GetString("kafka.consumer_group"),
		cfg)
	if err != nil {
		return nil, err
	}

	logger.Info("Successfully created kafka consumer")

	return consumerGroup, nil
}
