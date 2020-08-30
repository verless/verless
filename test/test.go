// Package test contains some simple utils for more readable tests.
// It is based on https://github.com/benbjohnson/testing.
package test

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Formatting strings for printing test results.
const (
	assertFormat    string = "\\033[31m%s:%d: %v\\033[39m\\n\\n"
	okFormat        string = "\u001B[31m%s:%d: unexpected error: %s\u001B[39m\n\n"
	equalsFormat    string = "\u001B[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\u001B[39m\n\n"
	notEqualsFormat string = "\u001B[31m%s:%d:\n\n\tnot exp: %#v\n\n\tgot: %#v\u001B[39m\n\n"
)

// Assert fails the test if the condition is false.
// It returns true, if the test did not fail.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) bool {
	if condition {
		return true
	}

	_, file, line, _ := runtime.Caller(1)
	// msg has to be this way, else params "v" for it do not work.
	fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
	tb.Fail()
	return false

}

// Ok fails the test if an err is not nil.
// It returns true, if the test did not fail.
func Ok(tb testing.TB, err error) bool {
	if err == nil {
		return true
	}

	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
	tb.Error(err)
	return false
}

// Equals fails the test if exp is not equal to act.
// It returns true, if the test did not fail.
func Equals(tb testing.TB, exp, act interface{}, options ...cmp.Option) bool {
	if cmp.Equal(exp, act, options...) {
		return true
	}

	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
	tb.Fail()
	return false
}

// NotEquals fails the test if exp is equal to act.
// It returns true, if the test did not fail.
func NotEquals(tb testing.TB, exp, act interface{}, options ...cmp.Option) bool {
	if !cmp.Equal(exp, act, options...) {
		return true
	}

	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("\033[31m%s:%d:\n\n\tnot exp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
	tb.Fail()
	return false
}

// ExpectError fails if the error is not the given error, using errors.Is.
// It returns true, if the test did not fail.
func ExpectedError(tb testing.TB, exp, act error) bool {
	if errors.Is(act, exp) {
		return true
	}

	fmt.Printf("expected error \n%v, got \n%v", exp, act)
	tb.Fail()
	return false
}
