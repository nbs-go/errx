package errx

func NewBuilder(namespace string, args ...SetOptionFn) *Builder {
	b := &Builder{
		errMap:    make(map[string]*Error),
		namespace: namespace,
	}

	// Evaluate options
	o := evaluateOptions(args)

	// Set fallback error and override namespace
	if o.fallbackErr != nil {
		b.fallbackErr = o.fallbackErr.Copy(WithNamespace(b.namespace))
	} else {
		b.fallbackErr = InternalError().Copy(WithNamespace(b.namespace))
	}

	return b
}

// Builder is an error builder with template namespace. All error produced will have a namespace
type Builder struct {
	errMap      map[string]*Error
	namespace   string
	fallbackErr *Error
}

// NewError create new error and ensure error is unique by it's code
func (b *Builder) NewError(code string, message string, args ...SetOptionFn) *Error {
	// Prevent setting error code that is equal with fallback error
	if code == b.fallbackErr.Code() {
		panic(DuplicateFallbackError)
	}

	// If args is empty, then set namespace
	if len(args) == 0 {
		args = []SetOptionFn{WithNamespace(b.namespace)}
	} else {
		// Else, merge arguments and override namespace
		args = append(args, WithNamespace(b.namespace))
	}

	// Create error
	err := NewError(code, message, args...)

	// Register error to dictionary, always overwrite existing
	b.errMap[err.Code()] = err

	return err
}

// Get retrieve error by Code, if not exist then return fallback error
func (b *Builder) Get(code string) *Error {
	err, ok := b.errMap[code]
	if ok {
		return err
	}
	return b.fallbackErr
}

// Namespace is getter function to retrieve builder namespace
func (b *Builder) Namespace() string {
	return b.namespace
}

func (b *Builder) FallbackError() *Error {
	return b.fallbackErr
}
