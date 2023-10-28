package minio

import (
	"context"

	"github.com/codfrm/cago/pkg/oss/oss"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
}

type Client struct {
	client *minio.Client
	core   *minio.Core
}

func New(cfg *Config) (oss.Client, error) {
	// Initialize minio core object.
	minioCore, err := minio.NewCore(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return &Client{client: minioClient, core: minioCore}, nil
}

func (c *Client) ListBuckets(ctx context.Context) ([]*oss.BucketInfo, error) {
	list, err := c.core.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]*oss.BucketInfo, 0, len(list))
	for _, v := range list {
		ret = append(ret, &oss.BucketInfo{Name: v.Name})
	}
	return ret, nil
}

func (c *Client) Bucket(ctx context.Context, bucket string) (oss.Bucket, error) {
	return &Bucket{
		client: c,
		bucket: bucket,
	}, nil
}
