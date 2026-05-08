auth:
	go run cmd/auth/main.go

blog:
	go run cmd/blog/main.go

notification:
	go run cmd/notification/main.go

kafka-up:
	docker compose up -d kafka

kafka-down:
	docker compose down