package errlog_test

import (
	"errors"
	"fmt"
	"github.com/bukodi/playground/errlog"
	"log/slog"
	"os"
	"testing"
)

func testFn() error {
	return errlog.WithStack(fmt.Errorf("error in testFn"))
}

func TestStacktrace(t *testing.T) {
	err1 := testFn()
	err2 := errlog.WithStack(err1)
	t.Logf("Err2 :%#v", err2)
}

func TestStacktraceMulti(t *testing.T) {
	err1a := testFn()
	err1b := testFn()
	err1 := errors.Join(err1a, err1b)
	err2 := errlog.WithStack(err1)
	t.Logf("Err2 :%#v", err2)

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	})))
	slog.Error("Error occured", errlog.ErrorKey, err2)
}
