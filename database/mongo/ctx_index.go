package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CtxIndex struct {
	ctx   context.Context
	index mongo.IndexView
}

// NewCtxIndex creates a new CtxIndex
// 对原本的 mongo.IndexView 进行了封装，这样在调用时就不需要传入 context.Context
func NewCtxIndex(ctx context.Context, index mongo.IndexView) *CtxIndex {
	return &CtxIndex{
		ctx:   ctx,
		index: index,
	}
}

func (c *CtxIndex) CreateOne(model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error) {
	return c.index.CreateOne(c.ctx, model, opts...)
}

func (c *CtxIndex) CreateMany(models []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	return c.index.CreateMany(c.ctx, models, opts...)
}

func (c *CtxIndex) DropOne(name string, opts ...*options.DropIndexesOptions) (bson.Raw, error) {
	return c.index.DropOne(c.ctx, name, opts...)
}

func (c *CtxIndex) DropAll(opts ...*options.DropIndexesOptions) (bson.Raw, error) {
	return c.index.DropAll(c.ctx, opts...)
}

func (c *CtxIndex) List(opts ...*options.ListIndexesOptions) (*mongo.Cursor, error) {
	return c.index.List(c.ctx, opts...)
}

func (c *CtxIndex) ListSpecifications(opts ...*options.ListIndexesOptions) ([]*mongo.IndexSpecification, error) {
	return c.index.ListSpecifications(c.ctx, opts...)
}
