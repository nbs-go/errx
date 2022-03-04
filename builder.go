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
	// Create error
	args = b.mergeArgs(args)
	err := NewError(code, message, args...)

	// Register error to dictionary, always overwrite existing
	b.registerError(err)

	return err
}

// Get retrieve error by Code, if no t exist then return fallback error
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

// FallbackError is getter function to retrieve FallbackError value
func (b *Builder) FallbackError() *Error {
	return b.fallbackErr
}

// CopyError take error input and override the namespace
func (b *Builder) CopyError(err *Error, args ...SetOptionFn) *Error {
	// Copy error and override namespace
	args = b.mergeArgs(args)
	bErr := err.Copy(args...)

	// Register error to map
	b.registerError(bErr)

	return bErr
}

func (b *Builder) mergeArgs(args []SetOptionFn) []SetOptionFn {
	if len(args) == 0 {
		return []SetOptionFn{WithNamespace(b.namespace)}
	}
	return append(args, WithNamespace(b.namespace))
}

func (b *Builder) registerError(err *Error) {
	// Check code not to collide with Fallback
	if b.fallbackErr.Code() == err.Code() {
		panic(DuplicateFallbackError)
	}

	b.errMap[err.Code()] = err
}
