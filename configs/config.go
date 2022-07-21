package configs

type Config struct {
	AppName string
	source  Source
	config  map[string]interface{}
}

func NewConfig(appName string, source Source, opt ...Option) (*Config, error) {
	options := &Options{}
	for _, opt := range opt {
		opt(options)
	}
	c := &Config{
		AppName: appName,
		source:  source,
	}
	return c, nil
}

func (c *Config) Scan(key string, value interface{}) error {
	return c.source.Scan(key, value)
}
