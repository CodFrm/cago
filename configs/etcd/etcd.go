package etcd

import (
	"context"
	"errors"
	"path"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/configs/file"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Config struct {
	Endpoints []string `env:"ETCD_ENDPOINT"`
	Username  string   `env:"ETCD_USERNAME"`
	Password  string   `env:"ETCD_PASSWORD"`
	Prefix    string   `env:"ETCD_PREFIX"`
}

type etcd struct {
	*clientv3.Client
	prefix        string
	serialization file.Serialization
}

func NewSource(cfg *Config, serialization file.Serialization) (configs.Source, error) {
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		Username:    cfg.Username,
		Password:    cfg.Password,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &etcd{
		Client:        cli,
		prefix:        cfg.Prefix,
		serialization: serialization,
	}, nil
}

func (e *etcd) Scan(key string, value interface{}) error {
	resp, err := e.Client.Get(context.Background(), path.Join(e.prefix, key))
	if err != nil {
		return err
	}
	if len(resp.Kvs) == 0 {
		b, err := e.serialization.Marshal(value)
		if err != nil {
			return err
		}
		if _, err := e.Client.Put(context.Background(), path.Join(e.prefix, key), string(b)); err != nil {
			return err
		}
		return errors.New("not found")
	}
	return e.serialization.Unmarshal(resp.Kvs[0].Value, value)
}
