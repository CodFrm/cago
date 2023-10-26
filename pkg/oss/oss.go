package oss

import (
	"context"
	"github.com/codfrm/cago/configs"
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

	return nil
}

func Default() oss.Client {
	return defaultClient
}

func Bucket() oss.Bucket {
	return defaultBucket
}
