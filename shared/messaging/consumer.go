package messaging

import "github.com/IBM/sarama"

type Consumer struct {
	consumer *sarama.ConsumerGroup
}
