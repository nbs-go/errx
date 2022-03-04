package errx_test

import (
	"errors"
	"github.com/nbs-go/errx"
	"testing"
)

func TestBuilder(t *testing.T) {
	// Create builder
	b := errx.NewBuilder("myapp")

	// Create and register error to builder
	err := b.NewError("E_FMT_1", "Invalid format")

	// Check namespace
	if err.Namespace() != b.Namespace() {
		t.Errorf("unexpected err namespace created by builder. Namespace = %s", err.Namespace())
	}

	// Get error from builder
	err2 := b.Get("E_FMT_1")
	if !errors.Is(err2, err) {
		t.Errorf("unexpected err not registered in Builder. ActualError = %s", err2)
	}
}

func TestOverrideNamespace(t *testing.T) {
	// Create builder
	b := errx.NewBuilder("myapp")

	// Create and register error to builder with namespace options
	err := b.NewError("E_FMT_1", "Invalid format", errx.WithNamespace("other-app"))

	// Check namespace
	if err.Namespace() != b.Namespace() {
		t.Errorf("unexpected err namespace created by builder. Namespace = %s", err.Namespace())
	}
}

func TestCustomFallbackError(t *testing.T) {
	// Create builder
	b := errx.NewBuilder("myapp",
		errx.FallbackError(errx.NewError("500", "Internal Server Error")))

	// Create and register error to builder with namespace options
	err := b.Get("E_FMT_1")

	// Check namespace
	if !errors.Is(err, b.FallbackError()) {
		t.Errorf("unexpected fallback error. Actual = %s, Expected = %s", err, b.FallbackError())
	}
}

func TestDuplicateFallbackError(t *testing.T) {
	b := errx.NewBuilder("myapp", errx.FallbackError(errx.NewError("500", "Internal Error")))
	defer RecoverPanic(t, errx.DuplicateFallbackError)()
	_ = b.NewError("500", "Another Error")
}

func TestBuilderCopyError(t *testing.T) {
	b := errx.NewBuilder("myapp1")
	srcErr := errx.NewError("ERR_1", "Resource not found", errx.WithNamespace("myapp2"))

	cpErr := b.CopyError(srcErr)

	if cpErr.Namespace() == srcErr.Namespace() || cpErr.Namespace() != "myapp1" {
		t.Errorf("unexpected namespace equal. CopiedNamespace = %s", cpErr.Namespace())
	}
}

func TestBuilderCopyErrorWithMetadata(t *testing.T) {
	b := errx.NewBuilder("myapp1")
	srcErr := errx.NewError("ERR_1", "Resource not found", errx.WithNamespace("myapp2"))

	cpErr := b.CopyError(srcErr, errx.WithMetadata(map[string]interface{}{
		"httpStatus": 400,
	}))

	meta := cpErr.Metadata()
	if len(meta) != 1 {
		t.Errorf("unexpected metadata in copied error. Metadata = %+v", cpErr.Metadata())
		return
	}

	v, ok := meta["httpStatus"]
	if !ok {
		t.Errorf("unexpected httpStatus not found in metadata. Metadata = %+v", cpErr.Metadata())
	}

	if v != 400 {
		t.Errorf("unexpected httpStatus in metadata. httpStatus = %d", v)
		return
	}
}

func RecoverPanic(t *testing.T, expectation error) func() {
	return func() {
		r := recover()
		if r == nil {
			t.Errorf("unexpected code did not panic")
			return
		}

		switch actual := r.(type) {
		case error:
			if !errors.Is(actual, expectation) {
				t.Errorf("unexpected error. got = %s", actual)
			}
		default:
			t.Errorf("unexpected type recovering from panic. r = %+v", actual)
		}

	}
}
