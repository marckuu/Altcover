package messaging

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewProducer(brokerAddresses []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()

	config.Version = sarama.V3_6_0_0
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Retry.Backoff = 100 * time.Millisecond

	p, err := sarama.NewSyncProducer(brokerAddresses, config)
	if err != nil {
		return &Producer{}, fmt.Errorf("ошибка при создании продюсера: %w", err)
	}

	return &Producer{
		producer: p,
		topic:    topic,
	}, nil
}

func (p *Producer) Produce(message []byte) error {
	prodMsg := &sarama.ProducerMessage{
		Topic:     p.topic,
		Key:       nil,
		Value:     sarama.ByteEncoder(message),
		Headers:   nil,
		Metadata:  nil,
		Offset:    0,
		Partition: 0,
		Timestamp: time.Time{},
	}

	if _, _, err := p.producer.SendMessage(prodMsg); err != nil {
		return fmt.Errorf("ошибка при отправке сообщения: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	if err := p.producer.Close(); err != nil {
		return fmt.Errorf("ошибка при закрытии продюсера: %w", err)
	}

	return nil
}
