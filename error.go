package errx

import (
	"errors"
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
	isSource  bool
}

// Error implement standard go error interface. If source error is exists then it will print error cause
func (e *Error) Error() string {
	errMsg := e.baseError()

	if len(e.traces) > 0 {
		errMsg += "\n  Traces => " + strings.Join(e.traces, "\n            ")
	}

	if e.sourceErr != nil {
		// Append CausedBy and traces
		errMsg += "\n  CausedBy => " + e.sourceErr.Error()
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

// Message is getter function to retrieve message value
func (e *Error) Message() string {
	return e.message
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

func (e *Error) wrapAndTrace(srcErr error) (*Error, []string) {
	// If srcErr is empty, then copy current error and its traces
	if srcErr == nil {
		return e.Copy(), copyTraces(e.traces)
	}

	// If srcErr error is equal to current error, Ignore source, copy current error and get traces from srcErr error
	if errors.Is(srcErr, e) {
		// Copy existing error and get traces from srcErr error
		var traces []string
		var sErr *Error
		ok := errors.As(srcErr, &sErr)
		if ok {
			traces = copyTraces(sErr.traces)
		}
		return e.Copy(), traces
	}

	// Init traces
	traces := make([]string, 0)

	// If srcErr error is a *errx.Error, then wrap error and move traces to current error
	if sErr, ok := srcErr.(*Error); ok && len(sErr.traces) > 0 {
		// Copy traces
		traces = copyTraces(sErr.traces)
		// Remove traces from srcErr error
		sErr.traces = nil
		// Set error as source
		sErr.isSource = true
	}

	return e.Wrap(srcErr), traces
}

func (e *Error) Trace(args ...SetOptionFn) *Error {
	// Get options
	o := evaluateOptions(args)

	// Wrap and trace error
	nErr, traces := e.wrapAndTrace(o.sourceErr)

	// Get trace
	ct := trace(o.skipTrace)
	nErr.traces = []string{ct}

	// If traces is exists, then merge
	if len(traces) > 0 {
		nErr.traces = append(nErr.traces, traces...)
	}

	// Merge metadata
	if len(o.metadata) > 0 {
		for k, v := range o.metadata {
			nErr.metadata[k] = v
		}
	}

	return nErr
}

// AddMetadata copy existing error and set new metadata
func (e *Error) AddMetadata(key string, value interface{}) *Error {
	// Copy error
	nErr := e.Copy()

	// Add metadata
	nErr.metadata[key] = value

	return nErr
}

// baseError print base error message with its codes
func (e *Error) baseError() string {
	if e.namespace == "" {
		if e.isSource {
			return fmt.Sprintf("[%s] %s", e.code, e.message)
		}
		return e.message
	}
	return fmt.Sprintf("%s: [%s] %s", e.namespace, e.code, e.message)
}
