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

POSTGRES_CONTAINER=microsvcdb
POSTGRES_USER=postgres
POSTGRES_DB=postgres

migrate-up:
	docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < ./migrations/000001_create_users_table.up.sql
	docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < ./migrations/000002_create_posts_table.up.sql

 migrate-down:
	docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < ./migrations/000001_create_users_table.down.sql
	docker exec -i $(POSTGRES_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < ./migrations/000002_create_posts_table.down.sql

