package errlog

import (
	"log/slog"
	"runtime"
	"strings"
)

const ErrorKey = "err"
const PackageKey = "pkg"
const AttrModule = "module"

func ErrorAttr(err error) slog.Attr {
	return slog.Any(ErrorKey, err)
}

func NewPkgLogger(parentLogger *slog.Logger) *slog.Logger {
	pkgName := "main"
	if pc, _, _, ok := runtime.Caller(1); ok {
		frames := runtime.CallersFrames([]uintptr{pc})
		frame, _ := frames.Next()
		callerFnName := frame.Func.Name()
		parts1 := strings.Split(callerFnName, "/")
		if len(parts1) >= 2 {
			parts2 := strings.Split(parts1[len(parts1)-1], ".")
			if len(parts2) == 2 {
				pkgName = parts2[0]
			}
		}
	}

	if parentLogger == nil {
		parentLogger = slog.Default()
	}
	return parentLogger.With("pkg", pkgName)
}
