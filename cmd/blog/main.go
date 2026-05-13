package main

import (
	"context"
	"log"

	blog "micro-blog/internal/blog"
)

func main() {
	ctx := context.Background()

	if err := blog.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
