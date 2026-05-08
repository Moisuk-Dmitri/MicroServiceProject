package notificationkafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type UserCreatedHandler interface {
	HandleUserCreated(ctx context.Context, event []byte) error
}

type Consumer struct {
	reader  *kafka.Reader
	handler UserCreatedHandler
}

func NewConsumer(brokers []string, topic string, groupID string, handler UserCreatedHandler) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
		handler: handler,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	log.Printf("notification kafka consumer started")
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			log.Printf("failed to read kafka message: %v", err)
			continue
		}

		log.Printf("received event: key=%s value=%s", string(msg.Key), string(msg.Value))
		if err := c.handler.HandleUserCreated(ctx, msg.Value); err != nil {
			log.Printf("failed to handle user.created event: %v", err)
			continue
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
