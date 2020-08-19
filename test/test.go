// Package test contains some simple utils for more readable tests.
// It is based on https://github.com/benbjohnson/testing.
package test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Formatting strings for printing test results.
const (
	assertFormat string = "\\033[31m%s:%d: %v\\033[39m\\n\\n"
	okFormat     string = "\u001B[31m%s:%d: unexpected error: %s\u001B[39m\n\n"
	equalsFormat string = "\u001B[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\u001B[39m\n\n"
)

// Assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(assertFormat, append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// Ok fails the test if an err is not nil.
func Ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(okFormat, filepath.Base(file), line, err.Error())
		tb.Error(err)
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}, options ...cmp.Option) {
	if !cmp.Equal(exp, act, options...) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(equalsFormat, filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
