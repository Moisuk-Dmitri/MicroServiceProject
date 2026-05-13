package authkafka

import (
	"context"
	"encoding/json"
	"fmt"
	"micro-blog/internal/events"

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

func (p *Producer) UserCreatedEvent(ctx context.Context, event events.UserCreatedEvent) error {
	eventJson, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("user created event: %w", err)
	}

	if err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.ID),
		Value: eventJson,
	}); err != nil {
		return fmt.Errorf("send user created event: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
