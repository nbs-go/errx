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
	xErr, _ := err.(*errx.Error)
	traces := xErr.Traces()

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

func nestedErr1() *errx.Error {
	return errx.NewError("ERR_1", "Bad Request").Trace(errx.Source(fmt.Errorf("invalid email format")))
}

func nestedErr2() error {
	return errx.Trace(nestedErr1())
}

func nestedErr3() error {
	return errx.Trace(nestedErr2())
}

func TestNestedTraceInternalError(t *testing.T) {
	err := nestedErr3()

	// Get traces
	xErr, _ := err.(*errx.Error)
	traces := xErr.Traces()

	expected := []string{
		"nbs-go/errx/error_test.go:99",
		"nbs-go/errx/error_test.go:95",
		"nbs-go/errx/error_test.go:91",
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

	t.Logf("Error = %s", err)
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

	if m := msgs[1]; !strings.HasSuffix(m, "nbs-go/errx/error_test.go:130") {
		t.Errorf("unexpected trace message. Trace = %s", m)
	}

	if m := msgs[2]; m != "  CausedBy => invalid phone format" {
		t.Errorf("unexpected cause. Caused By = %s", m)
	}

	t.Logf("Error = %s", err)
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

func TestAddMetadataOption(t *testing.T) {
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

	if msg := traces[0]; !strings.HasSuffix(msg, "nbs-go/errx/error_test.go:237") {
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

func TestAddMetadata(t *testing.T) {
	err1 := errx.NewError("ERR_1", "Resource not found", errx.AddMetadata("httpStatus", 404))
	err2 := err1.AddMetadata("overrideMessage", "Invoice not found")

	// Check metadata on err
	if meta := err1.Metadata(); len(meta) != 1 {
		t.Errorf("unexpected metadata on err1. Metadata = %+v", meta)
	}

	// Check metadata on err2
	if meta := err2.Metadata(); len(meta) != 2 {
		t.Errorf("unexpected metadata on err2. Metadata = %+v", meta)
	}
}

func newNestedError1() error {
	return errx.NewError("ERR_1", "customer.email is required").Trace()
}

func newNestedError2() error {
	return errx.NewError("ERR_2", "Failed to create customer").Trace(errx.Source(newNestedError1()))
}

func newNestedError3() error {
	return errx.NewError("400", "Bad Request").Trace(errx.Source(newNestedError2()))
}

func newNestedError4() error {
	return errx.Trace(newNestedError3())
}

func newNestedError5() error {
	return errx.Trace(newNestedError4())
}

func TestNestedTraceError(t *testing.T) {
	err := newNestedError5()

	xErr, ok := err.(*errx.Error)
	if !ok {
		t.Errorf("unexpected type of err. Type = %t", err)
		return
	}

	if len(xErr.Traces()) != 5 {
		t.Errorf("unexpected traces length. Length = %d", len(xErr.Traces()))
		return
	}

	t.Logf("Error = %s", err)
}

func TestTraceWithMetadata(t *testing.T) {
	err := errx.NewError("E_SRC_1", "This is source error", errx.AddMetadata("key", "value"))

	metaKey := "hiddenMessage"
	metaValue := "this is hidden message"
	err = err.Trace(errx.AddMetadata(metaKey, metaValue))

	m := err.Metadata()

	v, ok := m[metaKey]
	if !ok {
		t.Errorf("unexpected %s not found in metadata", metaKey)
		return
	}

	if v != metaValue {
		t.Errorf("unexpected metadata value = %s", metaValue)
	}
}

func BenchmarkNested(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = newNestedError5()
	}
}
