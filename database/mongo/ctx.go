package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CtxMongoDatabase struct {
	ctx      context.Context
	database *mongo.Database
}

// NewCtxMongoDatabase creates a new CtxMongoDatabase
// 对原本的 mongo.Database 进行了封装，这样在调用时就不需要传入 context.Context
func NewCtxMongoDatabase(ctx context.Context, database *mongo.Database) *CtxMongoDatabase {
	return &CtxMongoDatabase{
		ctx:      ctx,
		database: database,
	}
}

func (c *CtxMongoDatabase) Collection(name string, opts ...*options.CollectionOptions) *CtxCollection {
	return &CtxCollection{c.ctx, c.database.Collection(name, opts...)}
}

type CtxCollection struct {
	ctx        context.Context
	collection *mongo.Collection
}

func (c *CtxCollection) FindOne(filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return c.collection.FindOne(c.ctx, filter, opts...)
}

func (c *CtxCollection) Find(filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return c.collection.Find(c.ctx, filter, opts...)
}

func (c *CtxCollection) FindOneAndUpdate(filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	return c.collection.FindOneAndUpdate(c.ctx, filter, update, opts...)
}

func (c *CtxCollection) FindOneAndDelete(filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	return c.collection.FindOneAndDelete(c.ctx, filter, opts...)
}

func (c *CtxCollection) FindOneAndReplace(filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {
	return c.collection.FindOneAndReplace(c.ctx, filter, replacement, opts...)
}

func (c *CtxCollection) InsertOne(document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return c.collection.InsertOne(c.ctx, document, opts...)
}

func (c *CtxCollection) InsertMany(documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return c.collection.InsertMany(c.ctx, documents, opts...)
}

func (c *CtxCollection) DeleteOne(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return c.collection.DeleteOne(c.ctx, filter, opts...)
}

func (c *CtxCollection) DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return c.collection.DeleteMany(c.ctx, filter, opts...)
}

func (c *CtxCollection) UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.collection.UpdateOne(c.ctx, filter, update, opts...)
}

func (c *CtxCollection) UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.collection.UpdateMany(c.ctx, filter, update, opts...)
}

func (c *CtxCollection) ReplaceOne(filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	return c.collection.ReplaceOne(c.ctx, filter, replacement, opts...)
}

func (c *CtxCollection) CountDocuments(filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return c.collection.CountDocuments(c.ctx, filter, opts...)
}

func (c *CtxCollection) EstimatedDocumentCount(opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	return c.collection.EstimatedDocumentCount(c.ctx, opts...)
}

func (c *CtxCollection) Distinct(fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
	return c.collection.Distinct(c.ctx, fieldName, filter, opts...)
}

func (c *CtxCollection) Aggregate(pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	return c.collection.Aggregate(c.ctx, pipeline, opts...)
}

func (c *CtxCollection) Watch(pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	return c.collection.Watch(c.ctx, pipeline, opts...)
}

func (c *CtxCollection) Name() string {
	return c.collection.Name()
}

func (c *CtxCollection) Drop() error {
	return c.collection.Drop(c.ctx)
}

func (c *CtxCollection) BulkWrite(requests []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	return c.collection.BulkWrite(c.ctx, requests, opts...)
}

func (c *CtxCollection) Clone(opts ...*options.CollectionOptions) (*CtxCollection, error) {
	collection, err := c.collection.Clone(opts...)
	if err != nil {
		return nil, err
	}
	return &CtxCollection{c.ctx, collection}, nil
}

func (c *CtxCollection) Indexes() *CtxIndex {
	return NewCtxIndex(c.ctx, c.collection.Indexes())
}
