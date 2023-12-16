package permission

type Options struct {
	storageMatcher Matcher
}

type Option func(*Options)

func WithStorageMatcher(matcher Matcher) Option {
	return func(options *Options) {
		options.storageMatcher = matcher
	}
}

func WithFuncStorageMatcher(matcher FuncMatcher) Option {
	return func(options *Options) {
		options.storageMatcher = matcher
	}
}

type CheckOptions struct {
	storageMatcher Matcher
}

type CheckOption func(*CheckOptions)

func WithCheckStorageMatcher(matcher Matcher) CheckOption {
	return func(options *CheckOptions) {
		options.storageMatcher = matcher
	}
}

func WithCheckFuncStorageMatcher(matcher FuncMatcher) CheckOption {
	return func(options *CheckOptions) {
		options.storageMatcher = matcher
	}
}
