package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	blog "micro-blog/internal/blog"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	if err := blog.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
