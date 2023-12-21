package permission

type Options struct {
	matcher Matcher
}

type Option func(*Options)

func WithMatcher(matcher Matcher) Option {
	return func(options *Options) {
		options.matcher = matcher
	}
}

func WithFuncMatcher(matcher FuncMatcher) Option {
	return func(options *Options) {
		options.matcher = matcher
	}
}

type CheckOptions struct {
	matcher Matcher
}

type CheckOption func(*CheckOptions)

func WithCheckMatcher(matcher Matcher) CheckOption {
	return func(options *CheckOptions) {
		options.matcher = matcher
	}
}

func WithCheckFuncMatcher(matcher FuncMatcher) CheckOption {
	return func(options *CheckOptions) {
		options.matcher = matcher
	}
}
