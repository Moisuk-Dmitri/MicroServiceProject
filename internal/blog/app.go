package blog

import (
	"context"
	bloggrpcclient "micro-blog/internal/blog/grpcclient"
	bloghttp "micro-blog/internal/blog/http"
	blogservice "micro-blog/internal/blog/service"
)

func Run(ctx context.Context) error {
	authClient, err := bloggrpcclient.NewAuthClient("localhost:8081")
	if err != nil {
		return err
	}
	defer authClient.Close()

	service := blogservice.NewService()
	httpServer := bloghttp.NewServer(
		":8082",
		service,
		authClient,
	)

	return httpServer.Run()
}
