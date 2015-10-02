package selfcompile

import (
	"errors"
	"testing"
)

func TestErrors(t *testing.T) {
	var nilErr error
	if want, got := nilErr, combineErrors(); want != got {
		t.Errorf("got %q; want %q", got, want)
	}

	err := errors.New("hello")
	if want, got := err, combineErrors(err); want != got {
		t.Errorf("got %q; want %q", got, want)
	}

	if want, got := "2 errors: hello; hello", combineErrors(err, err).Error(); want != got {
		t.Errorf("got %q; want %q", got, want)
	}
}
