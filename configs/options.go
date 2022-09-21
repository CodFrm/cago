package configs

import "github.com/codfrm/cago/configs/file"

type Option func(*Options)

type Options struct {
	file          string
	serialization file.Serialization
}

func WithConfigFile(file string) Option {
	return func(options *Options) {
		options.file = file
	}
}

func WithSerialization(serialization file.Serialization) Option {
	return func(options *Options) {
		options.serialization = serialization
	}
}
