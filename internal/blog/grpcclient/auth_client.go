package bloggrpcclient

import (
	"context"
	authpb "micro-blog/proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client authpb.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: authpb.NewAuthServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *Client) ValidateToken(ctx context.Context, token string) (*authpb.ValidateTokenResponse, error) {
	resp, err := c.client.ValidateToken(ctx, &authpb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return &authpb.ValidateTokenResponse{}, err
	}

	return resp, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
