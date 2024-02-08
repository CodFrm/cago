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

	URL() string
}

type Bucket interface {
	PutObject(ctx context.Context, objectName string, r io.Reader, objectSize int64) error
	PreSignedPutObject(ctx context.Context, objectName string, expires time.Duration) (u *url.URL, err error)

	GetObject(ctx context.Context, objectName string) (Object, error)
	PreSignedGetObject(ctx context.Context, objectName string, expires time.Duration) (u *url.URL, err error)

	RemoveObject(ctx context.Context, objectName string) error

	FileURL(ctx context.Context, objectName string) (string, error)
}

type ObjectInfo struct {
	Key          string
	ETag         string
	Size         int64
	LastModified time.Time
	ContentType  string
}

type Object interface {
	io.ReadCloser
	Stat() (ObjectInfo, error)
}

type RespError interface {
	error
	StatusCode() int
}
