package config

import "os"

type Config struct {
	AuthHTTPPort     string
	BlogHTTPPort     string
	GRPCPort         string
	PostgresDSN      string
	KafkaAddr        string
	UserCreatedTopic string
}

func Load() Config {
	return Config{
		AuthHTTPPort: getEnv("AUTH_HTTP_PORT", "8080"),
		BlogHTTPPort: getEnv("BLOG_HTTP_PORT", "8082"),
		GRPCPort:     getEnv("GRPC_PORT", "8081"),
		PostgresDSN: getEnv(
			"POSTGRES_DSN",
			"postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable",
		),
		KafkaAddr:        getEnv("KAFKA_ADDR", "localhost:9092"),
		UserCreatedTopic: getEnv("USER_CREATED_TOPIC", "user.created"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
