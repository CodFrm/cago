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
	Endpoint        string `yaml:"endpoint"`
	URL             string `yaml:"url"`
	AccessKeyID     string `yaml:"accessKeyID"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	UseSSL          bool   `yaml:"useSSL"`
	Type            Type   `yaml:"type"`
	Bucket          string `yaml:"bucket"`
}

var defaultClient oss.Client
var defaultBucket oss.Bucket

// OSS 对象存储, 支持minio
func OSS(ctx context.Context, config *configs.Config) error {
	cfg := &Config{}
	err := config.Scan(ctx, "oss", cfg)
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
			UseSSL:          cfg.UseSSL,
			URL:             cfg.URL,
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

func DefaultBucket() oss.Bucket {
	return defaultBucket
}
