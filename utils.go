package errx

import (
	"fmt"
	"runtime"
)

func InternalError() *Error {
	return NewError("ERROR", "Internal Error")
}

// Trace wrap and trace error. If error is not *errx.Error then it will be wrapped into InternalError
// Else, it will add stack trace to error
func Trace(err error) error {
	if err == nil {
		return nil
	}

	// Check error type
	tErr, ok := err.(*Error)
	if !ok {
		// Set as internal error
		tErr = InternalError()
	}

	return tErr.Trace(Source(err), SkipTrace(2))
}

func Wrap(err error) *Error {
	return InternalError().Wrap(err)
}

// trace returns where in file and line the function being called
func trace(skip int) string {
	_, file, line, _ := runtime.Caller(skip + 1)
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
