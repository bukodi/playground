package errlog

import (
	"errors"
	"fmt"
	"testing"
)

type MyWrapper struct {
	err error
}

func (mw MyWrapper) Error() string {
	return mw.err.Error()
}

func (mw MyWrapper) Unwrap() error {
	return mw.err
}

func TestMultiErrorWithWrap(t *testing.T) {
	err1 := fmt.Errorf("first error")
	err2 := fmt.Errorf("second error")
	errJ := fmt.Errorf("contains three errors: %w, %w, %w", err1, MyWrapper{err2}, ErrBase1)

	t.Logf("Joined error  : %v", errJ)
	t.Logf("Joined error +: %+v", errJ)
	t.Logf("Joined error #: %#v", errJ)

	t.Logf("joined error is Err1Base: %t", errors.Is(errJ, ErrBase1))

	if errs, is := errJ.(interface{ Unwrap() []error }); is {
		for i, err := range errs.Unwrap() {
			t.Logf("%d. : %v", i, err)
		}
	} else {
		t.Error("Not wrapped error")
	}

}

func TestMultiErrorWithJoin(t *testing.T) {
	err1 := fmt.Errorf("first error")
	err2 := fmt.Errorf("second error")
	errJ := errors.Join(err1, err2, ErrBase1)

	t.Logf("Joined error  : %v", errJ)
	t.Logf("Joined error +: %+v", errJ)
	t.Logf("Joined error #: %#v", errJ)

	t.Logf("joined error is Err1Base: %t", errors.Is(errJ, ErrBase1))

	if errs, is := errJ.(interface{ Unwrap() []error }); is {
		for i, err := range errs.Unwrap() {
			t.Logf("%d. : %v", i, err)
		}
	} else {
		t.Error("Not wrapped error")
	}
}

func TestDeferErr(t *testing.T) {
	fn := func() (retErr error) {
		defer func() { retErr = errors.Join(retErr, errors.New("error in defer")) }()
		return errors.New("error in return")
	}
	err := fn()
	t.Logf("%+v", err)
}
