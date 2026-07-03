package messaging

import "github.com/IBM/sarama"

type MessageHandler struct {
	handler func([]byte)
}

func (mh *MessageHandler) Setup(csg sarama.ConsumerGroupSession) error {
	return nil
}

func (mh *MessageHandler) Cleanup(csg sarama.ConsumerGroupSession) error {
	return nil
}

func (mh *MessageHandler) ConsumeClaim(csg sarama.ConsumerGroupSession, cgc sarama.ConsumerGroupClaim) error {
	for msg := range cgc.Messages() {
		mh.handler(msg.Value)

		csg.MarkMessage(msg, "")
	}

	return nil
}
