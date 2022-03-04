package errx_test

import (
	"errors"
	"fmt"
	"github.com/nbs-go/errx"
	"strings"
	"testing"
)

func TestNewError(t *testing.T) {
	var err error = errx.NewError("ERR_1", "Invalid input")

	if err.Error() != "Invalid input" {
		t.Errorf("unexpected error output. Error = %s", err)
	}
}

func TestNewErrorWithNamespace(t *testing.T) {
	var err error = errx.NewError("ERR_1", "Invalid input", errx.WithNamespace("myapp"))

	if err.Error() != "myapp: [ERR_1] Invalid input" {
		t.Errorf("unexpected error output. Error = %s", err)
	}
}

func TestCopyCustomNamespaceAndMetadata(t *testing.T) {
	httpInternalErr := errx.InternalError().Copy(
		errx.WithNamespace("myapp"),
		errx.WithMetadata(map[string]interface{}{
			"httpStatus": 500,
		}))

	// Test error message
	if errMsg := httpInternalErr.Error(); errMsg != "myapp: [ERROR] Internal Error" {
		t.Errorf("unexpected error output. Error = %s", errMsg)
	}

	// Test metadata
	meta := httpInternalErr.Metadata()

	if v, ok := meta["httpStatus"]; !ok {
		t.Errorf("unexpected metadata httpStatus is not found in copied Internal Error")
	} else {
		if v != 500 {
			t.Errorf("unexpected httpStatus in metadata. ActualValue = %d", v)
		}
	}
}

func TestCopyExistingMetadata(t *testing.T) {
	err := errx.NewError("ERR_1", "Invalid input",
		errx.WithMetadata(map[string]interface{}{
			"httpStatus": 400,
		}))
	cpErr := err.Copy()

	// Test metadata
	meta := cpErr.Metadata()

	if v, ok := meta["httpStatus"]; !ok {
		t.Errorf("unexpected metadata httpStatus is not found in copied error")
	} else {
		if v != 400 {
			t.Errorf("unexpected httpStatus in metadata. ActualValue = %d", v)
		}
	}
}

func TestTraceInternalError(t *testing.T) {
	gErr := fmt.Errorf("invalid email format")
	err := errx.Trace(gErr)

	// Get trace
	traces := err.Traces()

	if len(traces) == 0 {
		t.Errorf("unexpected empty trace")
		return
	}

	// Check trace message
	trace := traces[0]
	if !strings.HasSuffix(trace, "nbs-go/errx/error_test.go:72") {
		t.Errorf("unexpected traced line. Trace = %s", trace)
	}
}

func TestNestedTraceInternalError(t *testing.T) {
	err := nestedErr3()

	// Get traces
	traces := err.Traces()

	expected := []string{
		"nbs-go/errx/error_test.go:122",
		"nbs-go/errx/error_test.go:118",
		"nbs-go/errx/error_test.go:114",
	}
	if len(traces) != len(expected) {
		t.Errorf("unexpected trace length. Length = %d", len(traces))
		return
	}

	// Check trace message
	for i, trace := range traces {
		if !strings.HasSuffix(trace, expected[i]) {
			t.Errorf("unexpected traced line. Trace = %s", trace)
		}
	}
}

func nestedErr1() *errx.Error {
	return errx.Trace(fmt.Errorf("invalid email format"))
}

func nestedErr2() *errx.Error {
	return errx.InternalError().Trace(errx.Source(nestedErr1()))
}

func nestedErr3() *errx.Error {
	return errx.Trace(nestedErr2())
}

func TestPrintWithCause(t *testing.T) {
	err := errx.InternalError().Trace(errx.Source(fmt.Errorf("invalid phone format")))

	errMsg := err.Error()

	// Split message by new line
	msgs := strings.Split(errMsg, "\n")

	if len(msgs) != 3 {
		t.Errorf("unexpected messages. MessageLen = %d", len(msgs))
		return
	}

	if m := msgs[0]; m != "Internal Error" {
		t.Errorf("unexpected base message. Message = %s", m)
	}

	if m := msgs[1]; m != "  CausedBy => invalid phone format" {
		t.Errorf("unexpected cause. Caused By = %s", m)
	}

	if m := msgs[2]; !strings.HasSuffix(m, "nbs-go/errx/error_test.go:126") {
		t.Errorf("unexpected trace message. Trace = %s", m)
	}
}

func TestTraceNil(t *testing.T) {
	err := errx.Trace(nil)

	if err != nil {
		t.Errorf("unexpected error must be nil. Error = %s", err)
	}
}

func TestWrapNil(t *testing.T) {
	err := errx.InternalError().Wrap(nil)

	if err != nil {
		t.Errorf("unexpected error must be nil. Error = %s", err)
	}
}

func TestUnwrap(t *testing.T) {
	srcErr := fmt.Errorf("invalid phone format")
	err := errx.Wrap(srcErr)

	// Unwrap error
	if unErr := errors.Unwrap(err); srcErr != unErr {
		t.Errorf("unexpected unwrapped error. Error = %s", unErr)
	}
}

func TestIsError(t *testing.T) {
	// Check generic error
	expected := errx.NewError("ERR_1", "fullName is required")
	actual := errx.NewError("ERR_1", "fullName is required")

	if !errors.Is(actual, expected) {
		t.Errorf("unexpected errx.Error.\n  Actual = %s\n  Expected = %s", actual, expected)
	}
}

func TestIsWrappedError(t *testing.T) {
	expected := errx.NewError("ERR_1", "fullName is required")
	actual := errx.InternalError().Wrap(expected)

	if !errors.Is(actual, expected) {
		t.Errorf("unexpected wrapped error.\n  Actual = %s\n  Expected = %s", actual, expected)
	}
}

func TestIsWrappedFmtError(t *testing.T) {
	expected := errx.NewError("ERR_1", "fullName is required")
	actual := fmt.Errorf("internal error. %w", expected)

	if !errors.Is(actual, expected) {
		t.Errorf("unexpected check error by wrapped by fmt.Errorf.\n  Actual = %s\n  Expected = %s", actual, expected)
	}
}

func TestIsGenericError(t *testing.T) {
	expected := fmt.Errorf("fullName is required")
	actual := errx.NewError("ERR_1", "fullName is required")

	if errors.Is(actual, expected) {
		t.Errorf("unexpected generic error must not satisfy errx.Error.\n  Actual = %s\n  Expected = %s", actual, expected)
	}
}

func TestAddMetadata(t *testing.T) {
	err := errx.NewError("ERR_1", "Resource not found", errx.AddMetadata("httpStatus", 400))

	meta := err.Metadata()

	v, ok := meta["httpStatus"]
	if !ok {
		t.Errorf("unexpected httpStatus metadata is not set. Metadata = %+v", meta)
		return
	}

	if v != 400 {
		t.Errorf("unexpected httpStatus value in metadata. status = %d", v)
	}
}

func TestTraceEmpty(t *testing.T) {
	err := errx.InternalError().Trace()

	traces := err.Traces()

	if len(traces) != 1 {
		t.Errorf("unexpected traces length. Length = %d", len(traces))
		return
	}

	if msg := traces[0]; !strings.HasSuffix(msg, "nbs-go/errx/error_test.go:231") {
		t.Errorf("unexpected trace message. Trace = %s", msg)
	}
}

func TestTraceErrorf(t *testing.T) {
	err := errx.InternalError().Trace(errx.Errorf("unexpected value not found"))

	uErr := err.Unwrap()

	if uErr.Error() != "unexpected value not found" {
		t.Errorf("unexpected traced Errorf. Error = %s", uErr)
	}
}

func TestErrorMessage(t *testing.T) {
	err := errx.NewError("ERR_1", "malformed token")

	if err.Message() != "malformed token" {
		t.Errorf("unexpected error message. Error = %s", err)
	}
}

func BenchmarkNested(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := nestedErr3()

		// Get traces
		traces := err.Traces()

		expected := []string{
			"nbs-go/errx/error_test.go:122",
			"nbs-go/errx/error_test.go:118",
			"nbs-go/errx/error_test.go:114",
		}
		if len(traces) != len(expected) {
			b.Errorf("unexpected trace length. Length = %d", len(traces))
			return
		}

		// Check trace message
		for it, trace := range traces {
			if !strings.HasSuffix(trace, expected[it]) {
				b.Errorf("unexpected traced line. Trace = %s", trace)
			}
		}
	}
}
