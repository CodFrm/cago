//go:build min

package configs

// min版本的初始化,不会自动根据配置文件中的source字段选择配置源
func (c *Config) init() error {
	return nil
}
