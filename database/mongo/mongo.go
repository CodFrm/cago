package mongo

import (
	"context"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/pkg/trace"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

type Config struct {
	Uri      string `yaml:"uri"`
	Database string `json:"database"`
}

var defaultClient *mongo.Client
var defaultDatabase *mongo.Database

func Mongo(ctx context.Context, config *configs.Config) error {
	cfg := &Config{}
	if err := config.Scan("mongo", cfg); err != nil {
		return err
	}
	opts := options.Client()
	if trace.Default() != nil {
		opts.Monitor = otelmongo.NewMonitor(otelmongo.WithTracerProvider(trace.Default()))
	}
	opts.ApplyURI(cfg.Uri)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}
	defaultClient = client
	defaultDatabase = client.Database(cfg.Database)
	return nil
}

func Ctx(ctx context.Context) *CtxMongoDatabase {
	return &CtxMongoDatabase{
		ctx:      ctx,
		database: defaultDatabase,
	}
}