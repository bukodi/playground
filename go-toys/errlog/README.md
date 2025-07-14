Useful link:
https://betterstack.com/community/guides/logging/logging-in-go/


# Error handling technics

## Checking error with compact if

Usual form:
```go
err := foo()
if err != nul {
	return err
} 
```
Compact form:
```go
if err := foo(); err != nul {
	return err
} 
```


## slog initialization
Wrong:
```go
var pkgLogger = slog.With("pkg", "mypkgname")
func Foo() {
	pkgLogger.Debug("Foo called")
}
```
When a user of this module sets the defaultLogger, it will not be reflected.

Ok:
```go
var pkgLogger = slog.With("pkg", "mypkgname")

func Foo() {
pkgLogger.Debug("Foo called")
}
```

