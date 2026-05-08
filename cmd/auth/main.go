package main

import (
	"log"
	auth "micro-blog/internal/auth"
)

func main() {
	if err := auth.Run(); err != nil {
		log.Fatal(err)
	}
}
