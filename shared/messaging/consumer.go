package messaging

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer sarama.ConsumerGroup
	topics   []string
}

func NewConsumer(brokerAddresses []string, groupID string, topics []string) (*Consumer, error) {
	config := sarama.NewConfig()

	config.Version = sarama.V3_6_0_0
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	c, err := sarama.NewConsumerGroup(brokerAddresses, groupID, config)
	if err != nil {
		return &Consumer{}, fmt.Errorf("ошибка при создании консюмера: %w", err)
	}

	return &Consumer{
		consumer: c,
		topics:   topics,
	}, nil
}

func (c *Consumer) Consume(handler func([]byte)) error {
	ctx := context.Background()

	for {
		if err := c.consumer.Consume(ctx, c.topics, &MessageHandler{handler: handler}); err != nil {
			return fmt.Errorf("ошибка при обработке сообщения: %w", err)
		}
	}
}

func (c *Consumer) Close() error {
	if err := c.consumer.Close(); err != nil {
		return fmt.Errorf("ошибка при закрытии консюмера: %w", err)
	}

	return nil
}
