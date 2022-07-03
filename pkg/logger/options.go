package logger

type Option func(*Options)

type Options struct {
	level string
	debug bool
	loki  *LokiConfig
}

type LokiConfig struct {
	Level string
	Url   string
}

func WithLoki(config *LokiConfig) Option {
	return func(options *Options) {
		options.loki = config
	}
}

func Level(level string) Option {
	return func(options *Options) {
		options.level = level
	}
}

func Debug() Option {
	return func(options *Options) {
		options.debug = true
	}
}
