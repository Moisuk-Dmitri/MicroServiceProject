package notification

import (
	"context"
	notificationkafka "micro-blog/internal/notification/kafka"
	"micro-blog/internal/platform/config"
)

func Run(ctx context.Context) error {
	cfg := config.Load()

	brokers := []string{cfg.KafkaAddr}
	topic := cfg.KafkaUserCreatedTopic
	groupID := cfg.KafkaGroupID
	handler := notificationkafka.NewHandler()

	consumer := notificationkafka.NewConsumer(brokers, topic, groupID, handler)
	defer consumer.Close()

	return consumer.Run(ctx)
}
