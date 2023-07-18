//go:build !min

package configs

import (
	"path"

	"github.com/codfrm/cago/configs/etcd"
)

// 非min版本的初始化,自动根据配置文件中的source字段选择配置源
func (c *Config) init() error {
	configSource := ""
	if err := c.source.Scan("source", &configSource); err != nil {
		return err
	}
	switch configSource {
	case "etcd":
		etcdConfig := &etcd.Config{}
		if err := c.source.Scan("etcd", etcdConfig); err != nil {
			return err
		}
		var err error
		etcdConfig.Prefix = path.Join(etcdConfig.Prefix, string(c.Env), c.AppName)
		c.source, err = etcd.NewSource(etcdConfig, c.serialization)
		if err != nil {
			return err
		}
	}
	return nil
}
