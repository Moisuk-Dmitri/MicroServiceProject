package blog

import (
	"context"
	"database/sql"
	bloggrpcclient "micro-blog/internal/blog/grpcclient"
	bloghttp "micro-blog/internal/blog/http"
	blogrepository "micro-blog/internal/blog/repository"
	blogservice "micro-blog/internal/blog/service"
	"micro-blog/internal/platform/config"
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

	return httpServer.Run()
}
