package main

import (
	"context"
	"log"
	auth "micro-blog/internal/auth"
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

	if err := auth.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
