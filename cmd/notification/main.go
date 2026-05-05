package main

import (
	"log"
	"time"
)

func main() {
	log.Printf("notification service started")

	for {
		log.Printf("notification running...")
		time.Sleep(5 * time.Second)
	}
}
