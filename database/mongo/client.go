package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Client struct {
	client   *mongo.Client
	database string
}

func (c *Client) Client() *mongo.Client {
	return c.client
}

func (c *Client) Database(ctx context.Context) *CtxMongoDatabase {
	return NewCtxMongoDatabase(ctx, c.client.Database(c.database))
}
