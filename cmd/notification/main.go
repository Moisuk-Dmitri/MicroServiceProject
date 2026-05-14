package main

import (
	"context"
	"log"
	"micro-blog/internal/notification"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	if err := notification.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
