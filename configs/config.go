package configs

import (
	"reflect"
	"strings"

	"github.com/codfrm/cago/configs/file"
	"github.com/codfrm/cago/configs/source"
)

type Env string

const (
	DEV  Env = "dev"
	TEST Env = "test"
	PRE  Env = "pre"
	PROD Env = "prod"
)

type Config struct {
	AppName       string
	Version       string
	Env           Env
	Debug         bool
	source        source.Source
	serialization file.Serialization
}

func NewConfig(appName string, opt ...Option) (*Config, error) {
	options := &Options{
		file:          "./configs/config.yaml",
		serialization: file.Yaml(),
	}
	for _, o := range opt {
		o(options)
	}
	var s source.Source
	if options.source == nil {
		var err error
		s, err = file.NewSource(options.file, options.serialization)
		if err != nil {
			return nil, err
		}
	} else {
		s = options.source
	}
	var env Env
	if err := s.Scan("env", &env); err != nil {
		return nil, err
	}
	var debug bool
	if err := s.Scan("debug", &debug); err != nil {
		return nil, err
	}
	version := ""
	if err := s.Scan("version", &version); err != nil {
		return nil, err
	}

	c := &Config{
		AppName:       appName,
		Debug:         debug,
		Env:           env,
		Version:       version,
		source:        s,
		serialization: options.serialization,
	}
	if err := c.init(); err != nil {
		return nil, err
	}
	defaultConfig = c
	return c, nil
}

func (c *Config) Scan(key string, value interface{}) error {
	return c.source.Scan(key, value)
}

func (c *Config) findKey(key string, value interface{}) error {
	keys := strings.Split(key, ".")
	if len(keys) == 1 {
		return c.source.Scan(key, value)
	}
	valueMap := make(map[string]interface{})
	if err := c.source.Scan(keys[0], &valueMap); err != nil {
		return err
	}
	for i := 1; i < len(keys); i++ {
		if v, ok := valueMap[keys[i]]; ok {
			if i == len(keys)-1 {
				reflect.ValueOf(value).Elem().Set(reflect.ValueOf(v))
				return nil
			}
			valueMap = v.(map[string]interface{})
		} else {
			return nil
		}
	}
	// 使用反射复制valueMap到value
	reflect.ValueOf(value).Elem().Set(reflect.ValueOf(valueMap))
	return nil
}

func (c *Config) String(key string) string {
	var str string
	if err := c.findKey(key, &str); err != nil {
		return ""
	}
	return str
}

func (c *Config) Has(key string) (bool, error) {
	return c.source.Has(key)
}
