package client

import (
	"context"
	pb "library-api-user/proto/book"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BookClient struct {
	client pb.BookServiceClient
}

func NewBookClient(addr string) (*BookClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &BookClient{client: pb.NewBookServiceClient(conn)}, nil
}

func (c *BookClient) DecreaseStock(ctx context.Context, bookID uint64) error {
	_, err := c.client.DecreaseStock(ctx, &pb.DecreaseStockRequest{BookId: bookID})
	return err
}

func (c *BookClient) IncreaseStock(ctx context.Context, bookID uint64) error {
	_, err := c.client.IncreaseStock(ctx, &pb.IncreaseStockRequest{BookId: bookID})
	return err
}
