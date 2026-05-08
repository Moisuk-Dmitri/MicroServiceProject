package notificationkafka

import (
	"context"
	"log"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HandleUserCreated(ctx context.Context, event []byte) error {
	log.Printf("handle user.created event: %s", string(event))
	return nil
}
