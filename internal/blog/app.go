package blog

import (
	"context"
	"database/sql"
	"fmt"
	bloggrpcclient "micro-blog/internal/blog/grpcclient"
	bloghttp "micro-blog/internal/blog/http"
	blogrepository "micro-blog/internal/blog/repository"
	blogservice "micro-blog/internal/blog/service"
	"micro-blog/internal/platform/config"
	"time"

	_ "github.com/lib/pq"
)

func Run(ctx context.Context) error {
	cfg := config.Load()
	if cfg.BlogHTTPPort[0] != ':' {
		cfg.BlogHTTPPort = ":" + cfg.BlogHTTPPort
	}
	if cfg.GRPCPort[0] != ':' {
		cfg.GRPCPort = ":" + cfg.GRPCPort
	}

	authClient, err := bloggrpcclient.NewAuthClient(cfg.GRPCPort)
	if err != nil {
		return err
	}
	defer authClient.Close()

	db, err := sql.Open("postgres", cfg.PostgresDSN)
	if err != nil {
		return err
	}
	if err := db.PingContext(ctx); err != nil {
		return err
	}
	repository := blogrepository.NewRepository(db)

	service := blogservice.NewService(repository)

	httpServer := bloghttp.NewServer(
		cfg.BlogHTTPPort,
		service,
		authClient,
	)

	errCh := make(chan error, 1)

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
		authClient.Close()

		return nil
	}
}
