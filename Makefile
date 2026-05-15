auth:
	go run cmd/auth/main.go

blog:
	go run cmd/blog/main.go

notification:
	go run cmd/notification/main.go

kafka-up:
	docker compose up -d kafka

postgres-up:
	docker compose up -d postgres

docker-up:
	docker compose up -d

docker-down:
	docker compose down

project-up:
	docker compose up --build

project-down:
	docker compose down