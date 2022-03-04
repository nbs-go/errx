package errx

import (
	"fmt"
	"strings"
)

// NewError initiates a new error instance
func NewError(code string, message string, args ...SetOptionFn) *Error {
	// Init error
	err := &Error{
		code:     code,
		message:  message,
		metadata: make(map[string]interface{}),
		traces:   make([]string, 0),
	}

	// Evaluate options
	o := evaluateOptions(args)

	// Override namespace if set in options
	if o.namespace != "" {
		err.namespace = o.namespace
	}

	// Set metadata value
	if len(o.metadata) > 0 {
		err.metadata = o.metadata
	}

	return err
}

// Error is an immutable object. print error meaningful message and stack trace for easier error tracing
type Error struct {
	code      string
	message   string
	namespace string
	metadata  map[string]interface{}
	sourceErr error
	traces    []string
}

// Error implement standard go error interface. If source error is exists then it will print error cause
func (e *Error) Error() string {
	errMsg := e.baseError()

	if e.sourceErr != nil {
		// Append CausedBy and traces
		errMsg += "\n  CausedBy => " + e.sourceErr.Error()
	}

	if len(e.traces) > 0 {
		errMsg += "\n  Traces => " + strings.Join(e.traces, "\n            ")
	}

	return errMsg
}

// Unwrap implements xerrors.Wrapper interface
func (e *Error) Unwrap() error {
	return e.sourceErr
}

// Copy duplicate error traces. Available options is WithNamespace, WithMetadata and CopySource
func (e *Error) Copy(args ...SetOptionFn) *Error {
	err := &Error{
		code:      e.code,
		message:   e.message,
		namespace: e.namespace,
		sourceErr: e.sourceErr,
		traces:    []string{},
	}

	o := evaluateOptions(args)

	// If namespace is set, then override namespace
	if o.namespace != "" {
		err.namespace = o.namespace
	}

	// If metadata is set, then override
	if len(o.metadata) > 0 {
		err.metadata = o.metadata
	} else {
		// Else, copy metadata
		err.metadata = copyMetadata(e.metadata)
	}

	return err
}

// Is implements function that will be called by errors.Is for error comparison.
// Actual error namespace and code value must equal with Expected ones
func (e *Error) Is(err error) bool {
	expected, ok := err.(*Error)
	if !ok {
		return false
	}

	return expected.Namespace() == e.Namespace() &&
		expected.Code() == e.Code()
}

// Code is getter function to retrieve error code value
func (e *Error) Code() string {
	return e.code
}

// Namespace is getter function to retrieve error namespace
func (e *Error) Namespace() string {
	return e.namespace
}

// Metadata is getter function to retrieve metadata value
func (e *Error) Metadata() map[string]interface{} {
	return e.metadata
}

// Traces is getter function to retrieve traces value
func (e *Error) Traces() []string {
	return e.traces
}

func (e *Error) Wrap(err error) *Error {
	if err == nil {
		return nil
	}

	// Copy error
	nErr := e.Copy()

	// Set source
	nErr.sourceErr = err

	return nErr
}

func (e *Error) Trace(args ...SetOptionFn) *Error {
	// Get options
	o := evaluateOptions(args)

	var nErr *Error
	if o.sourceErr != nil {
		nErr = e.Wrap(o.sourceErr)
	} else {
		nErr = e.Copy()
	}

	// Get trace
	ct := trace(o.skipTrace)
	nErr.traces = []string{ct}

	// Copy traces from errx.Error wrapper if exists
	if sErr, ok := nErr.sourceErr.(*Error); ok && len(sErr.traces) > 0 {
		traces := copyTraces(sErr.traces)
		nErr.traces = append(nErr.traces, traces...)
	}

	return nErr
}

// baseError print base error message with its codes
func (e *Error) baseError() string {
	if e.namespace == "" {
		return e.message
	}
	return fmt.Sprintf("%s: [%s] %s", e.namespace, e.code, e.message)
}
