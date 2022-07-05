package config

import (
	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AppName string
	config  map[string]interface{}
}

func NewConfig(appName string, source Source, opt ...Option) (*Config, error) {
	options := &Options{}
	for _, opt := range opt {
		opt(options)
	}
	b, err := source.Read()
	if err != nil {
		return nil, err
	}
	c := &Config{
		AppName: appName,
		config:  make(map[string]interface{}),
	}
	if err := yaml.Unmarshal(b, &c.config); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) Scan(key string, value interface{}) error {
	b, err := yaml.Marshal(c.config[key])
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(b, value); err != nil {
		return err
	}
	return env.Parse(value)
}
