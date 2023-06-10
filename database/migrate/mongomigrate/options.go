package mongomigrate

type Options struct {
	hasPre bool // 有预发布版本
}

type Option func(o *Options)

func HasPre() Option {
	return func(o *Options) {
		o.hasPre = true
	}
}
