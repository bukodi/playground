package errlog

import (
	"errors"
	"fmt"
	"testing"
)

var ErrBase1 = fmt.Errorf("base error 1")
var ErrBase2 = fmt.Errorf("base error 2")

type ErrType1 struct{}

func (et1 ErrType1) Error() string {
	return "typed error 1"
}

func (et1 ErrType1) String() string {
	return "typed error 1 as go string"
}
func (et1 ErrType1) GoString() string {
	return "typed error 1 as go string"
}

func TestFormat(t *testing.T) {
	err1 := ErrType1{}
	t.Logf("err1   :%v", err1)
	t.Logf("err1 + :%+v", err1)
	t.Logf("err1 # :%#v", err1)
	t.Logf("err1 s :%s", err1)

	wrapped2 := fmt.Errorf("caused by: %w", err1)
	t.Logf("wrapped2   :%v", wrapped2)
	t.Logf("wrapped2 + :%+v", wrapped2)
	t.Logf("wrapped2 # :%#v", wrapped2)
	t.Logf("wrapped2 s :%s", wrapped2)
}

func TestAsIs(t *testing.T) {
	err1 := fmt.Errorf("something happened (%w)", ErrBase1)

	t.Logf("err1 is ErrBase1: %t", errors.Is(err1, ErrBase1))

	wrapped2 := fmt.Errorf("caused by: %w", err1)
	t.Logf("wrapped2 is Err1Base: %t", errors.Is(wrapped2, ErrBase1))

	wrappedt2 := fmt.Errorf("caused by: %w", ErrType1{})
	wrappedt3 := fmt.Errorf("caused by: %w", wrappedt2)
	t.Logf("wrapped3  :%v", wrappedt3)
	t.Logf("wrapped3 + :%+v", wrappedt3)
	t.Logf("wrapped3 # :%#v", wrappedt3)
	var typedErr ErrType1
	errors.As(wrappedt3, &typedErr)
	t.Logf("wrapped3 as ErrType1: %v", typedErr)

}
