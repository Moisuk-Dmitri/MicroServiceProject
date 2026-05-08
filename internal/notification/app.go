package notification

import (
	"context"
	notificationkafka "micro-blog/internal/notification/kafka"
)

func Run(ctx context.Context) error {
	brokers := []string{"localhost:9092"}
	topic := "user.created"
	groupID := "notification-service"
	handler := notificationkafka.NewHandler()

	consumer := notificationkafka.NewConsumer(brokers, topic, groupID, handler)
	defer consumer.Close()

	return consumer.Run(ctx)
}
