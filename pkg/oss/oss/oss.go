package oss

import (
	"context"
	"io"
	"net/url"
	"time"
)

// BucketInfo container for bucket metadata.
type BucketInfo struct {
	// The name of the bucket.
	Name string `json:"name"`
}

type Client interface {
	ListBuckets(ctx context.Context) ([]*BucketInfo, error)

	Bucket(ctx context.Context, bucket string) (Bucket, error)
}

type Bucket interface {
	PutObject(ctx context.Context, objectName string, r io.Reader) error
	PreSignedPutObject(ctx context.Context, objectName string, expires time.Duration) (u *url.URL, err error)

	GetObject(ctx context.Context, objectName string) (io.ReadCloser, error)
	PreSignedGetObject(ctx context.Context, objectName string, expires time.Duration) (u *url.URL, err error)

	RemoveObject(ctx context.Context, objectName string) error
}
