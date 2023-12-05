package etcd

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/codfrm/cago/configs"

	"github.com/codfrm/cago/configs/file"
	"github.com/codfrm/cago/configs/source"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func init() {
	configs.RegistrySource("etcd", func(cfg *configs.Config, serialization file.Serialization) (source.Source, error) {
		etcdConfig := &Config{}
		if err := cfg.Scan("etcd", etcdConfig); err != nil {
			return nil, err
		}
		var err error
		etcdConfig.Prefix = path.Join(etcdConfig.Prefix, string(cfg.Env), cfg.AppName)
		s, err := NewSource(etcdConfig, serialization)
		if err != nil {
			return nil, err
		}
		return s, nil
	})
}

type Config struct {
	Endpoints []string
	Username  string
	Password  string
	Prefix    string
}

type etcd struct {
	*clientv3.Client
	prefix        string
	serialization file.Serialization
}

func NewSource(cfg *Config, serialization file.Serialization) (source.Source, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		Username:    cfg.Username,
		Password:    cfg.Password,
		DialTimeout: 10 * time.Second,
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
		return fmt.Errorf("etcd %w: %s", source.ErrNotFound, key)
	}
	return e.serialization.Unmarshal(resp.Kvs[0].Value, value)
}

func (e *etcd) Has(key string) (bool, error) {
	resp, err := e.Client.Get(context.Background(), path.Join(e.prefix, key))
	if err != nil {
		return false, err
	}
	return len(resp.Kvs) > 0, nil
}
