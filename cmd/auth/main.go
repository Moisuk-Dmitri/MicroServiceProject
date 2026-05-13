package main

import (
	"context"
	"log"
	auth "micro-blog/internal/auth"
)

func main() {
	ctx := context.Background()

	if err := auth.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
