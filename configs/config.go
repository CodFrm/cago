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
	TEST Env = "test"
	PROD Env = "prod"
)

type Config struct {
	AppName string
	Env     Env
	source  source2.Source
	config  map[string]interface{}
}

func NewConfig(appName string, opt ...Option) (*Config, error) {
	options := &Options{
		file:          "./configs/config.yaml",
		serialization: file.Yaml(),
	}
	for _, o := range opt {
		o(options)
	}
	source, err := file.NewSource(options.file, options.serialization)
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
		source, err = etcd.NewSource(etcdConfig, options.serialization)
		if err != nil {
			return nil, err
		}
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
