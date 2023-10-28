package oss

import (
	"context"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/pkg/oss/minio"
	"github.com/codfrm/cago/pkg/oss/oss"
)

type Type string

const (
	Minio Type = "minio"
)

type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	Type            Type
	Bucket          string
}

var defaultClient oss.Client
var defaultBucket oss.Bucket

// OSS 对象存储, 支持minio
func OSS(ctx context.Context, config *configs.Config) error {
	cfg := &Config{}
	err := config.Scan("oss", cfg)
	if err != nil {
		return err
	}
	var cli oss.Client
	switch cfg.Type {
	case Minio:
		cli, err = minio.New(&minio.Config{
			Endpoint:        cfg.Endpoint,
			AccessKeyID:     cfg.AccessKeyID,
			SecretAccessKey: cfg.SecretAccessKey,
		})
		if err != nil {
			return err
		}
	}
	defaultClient = newWrapClient(cli, newWrap())
	defaultBucket, err = cli.Bucket(ctx, cfg.Bucket)
	if err != nil {
		return err
	}
	return nil
}

func Default() oss.Client {
	return defaultClient
}

func Bucket() oss.Bucket {
	return defaultBucket
}
