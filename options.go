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

func AddMetadata(key string, value interface{}) SetOptionFn {
	return func(o *options) {
		o.metadata[key] = value
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

func Source(err error) SetOptionFn {
	return func(o *options) {
		o.sourceErr = err
	}
}

type options struct {
	namespace   string
	metadata    map[string]interface{}
	skipTrace   int
	fallbackErr *Error
	sourceErr   error
}

type SetOptionFn = func(*options)

func defaultOptions() *options {
	return &options{
		metadata:  make(map[string]interface{}),
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
