package errx

import (
	"fmt"
	"runtime"
)

func InternalError() *Error {
	return NewError("ERROR", "Internal Error")
}

// TraceError wrap and trace error. If error is not *errx.Error then it will be wrapped into InternalError
// Else, it will add stack trace to error
func TraceError(err error) *Error {
	// Check error type
	tErr, ok := err.(*Error)
	if !ok {
		// Set as internal error
		tErr = InternalError()
	}

	return tErr.Trace(err, SkipTrace(2))
}

func Wrap(err error) *Error {
	return InternalError().Wrap(err)
}

// trace returns where in file and line the function being called
func trace(skip int) string {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "<?>:<?>"
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func copyMetadata(m1 map[string]interface{}) map[string]interface{} {
	m2 := make(map[string]interface{})
	for k, v := range m1 {
		m2[k] = v
	}
	return m2
}

func copyTraces(m1 []string) []string {
	m2 := make([]string, len(m1))
	for k, v := range m1 {
		m2[k] = v
	}
	return m2
}
