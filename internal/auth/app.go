package auth

import (
	authgrpc "micro-blog/internal/auth/grpc"
	authhttp "micro-blog/internal/auth/http"
	authkafka "micro-blog/internal/auth/kafka"
	authservice "micro-blog/internal/auth/service"
)

func Run() error {
	producer := authkafka.NewProducer("localhost:9092", "user.created")
	defer producer.Close()

	service := authservice.NewService(producer)

	grpcServer := authgrpc.NewServer(":8081", service)
	httpServer := authhttp.NewServer(":8080", service)

	errCh := make(chan error, 2)

	go func() {
		errCh <- grpcServer.Run()
	}()

	go func() {
		errCh <- httpServer.Run()
	}()

	return <-errCh
}
