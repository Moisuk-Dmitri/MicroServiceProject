package authkafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(addr string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(addr),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) UserCreatedEvent(ctx context.Context) error {
	event := `{"user_id":"1"}`

	err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("user-1"),
		Value: []byte(event),
	})
	if err != nil {
		return fmt.Errorf("send user created event: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
