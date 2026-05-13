package auth

import (
	"context"
	"database/sql"
	"fmt"
	authgrpc "micro-blog/internal/auth/grpc"
	authhttp "micro-blog/internal/auth/http"
	authkafka "micro-blog/internal/auth/kafka"
	authrepository "micro-blog/internal/auth/repository"
	authservice "micro-blog/internal/auth/service"
	"micro-blog/internal/platform/config"
	"time"

	_ "github.com/lib/pq"
)

func Run(ctx context.Context) error {
	cfg := config.Load()

	db, err := sql.Open(
		"postgres",
		cfg.PostgresDSN,
	)
	if err != nil {
		return fmt.Errorf("open postgres: %w", err)
	}
	defer db.Close()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping postgres: %w", err)
	}

	repository := authrepository.NewUserRepository(db)

	producer := authkafka.NewProducer(cfg.KafkaAddr, cfg.UserCreatedTopic)
	defer producer.Close()

	service := authservice.NewService(producer, repository)

	grpcServer := authgrpc.NewServer(cfg.GRPCPort, service)
	httpServer := authhttp.NewServer(cfg.HTTPPort, service)

	errCh := make(chan error, 2)

	go func() {
		if err := grpcServer.Run(); err != nil {
			errCh <- err
		}
	}()

	go func() {
		if err := httpServer.Run(); err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("http server shutdown: %w", err)
		}

		return nil
	}
}
