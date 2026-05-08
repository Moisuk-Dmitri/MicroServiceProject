package main

import (
	"context"
	"log"
	"micro-blog/internal/notification"
)

func main() {
	ctx := context.Background()

	if err := notification.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
