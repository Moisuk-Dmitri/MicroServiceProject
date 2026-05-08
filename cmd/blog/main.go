package main

import (
	"log"

	blog "micro-blog/internal/blog"
)

func main() {
	if err := blog.Run(); err != nil {
		log.Fatal(err)
	}
}
