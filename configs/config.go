package configs

import (
	"path"

	"github.com/codfrm/cago/configs/etcd"
	"github.com/codfrm/cago/configs/file"
	source2 "github.com/codfrm/cago/configs/source"
)

type Env string

const (
	DEV  Env = "dev"
	PROD Env = "prod"
)

type Config struct {
	AppName string
	Env     Env
	source  source2.Source
	config  map[string]interface{}
}

func NewConfig(appName string, opt ...Option) (*Config, error) {
	source, err := file.NewSource("./configs/config.yaml", file.Yaml())
	if err != nil {
		return nil, err
	}
	configSource := ""
	if err := source.Scan("source", &configSource); err != nil {
		return nil, err
	}
	var env Env
	if err := source.Scan("env", &env); err != nil {
		return nil, err
	}

	switch configSource {
	case "etcd":
		etcdConfig := &etcd.Config{}
		if err := source.Scan("etcd", etcdConfig); err != nil {
			return nil, err
		}
		var err error
		etcdConfig.Prefix = path.Join(etcdConfig.Prefix, string(env), appName)
		source, err = etcd.NewSource(etcdConfig, file.Yaml())
		if err != nil {
			return nil, err
		}
	}

	options := &Options{}
	for _, opt := range opt {
		opt(options)
	}
	c := &Config{
		AppName: appName,
		Env:     env,
		source:  source,
	}
	return c, nil
}

func (c *Config) Scan(key string, value interface{}) error {
	return c.source.Scan(key, value)
}
