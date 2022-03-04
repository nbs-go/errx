package errx

func WithNamespace(namespace string) SetOptionFn {
	return func(o *options) {
		o.namespace = namespace
	}
}

func WithMetadata(metadata map[string]interface{}) SetOptionFn {
	return func(o *options) {
		o.metadata = metadata
	}
}

func SkipTrace(skip int) SetOptionFn {
	return func(o *options) {
		o.skipTrace = skip
	}
}

func FallbackError(err *Error) SetOptionFn {
	return func(o *options) {
		o.fallbackErr = err
	}
}

type options struct {
	namespace   string
	metadata    map[string]interface{}
	skipTrace   int
	fallbackErr *Error
}

type SetOptionFn = func(*options)

func defaultOptions() *options {
	return &options{
		skipTrace: 1,
	}
}

func evaluateOptions(args []SetOptionFn) *options {
	optCopy := defaultOptions()
	for _, fn := range args {
		fn(optCopy)
	}
	return optCopy
}
