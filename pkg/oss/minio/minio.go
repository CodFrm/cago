package minio

import (
	context "context"
	oss2 "github.com/codfrm/cago/pkg/oss"
	"github.com/codfrm/cago/pkg/oss/oss"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	client *minio.Client
	core   *minio.Core
}

func NewMinio(cfg *oss2.Config) (oss.Client, error) {
	endpoint := "192.168.1.136:9000"
	accessKeyID := "wCDPeNpuuxcpBU8l7oes"
	secretAccessKey := "fVyO7LQSvezIPPN9fL1uhRXaatxrpb7zU45eGFxm"

	// Initialize minio core object.
	minioCore, err := minio.NewCore(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
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
