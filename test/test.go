// Package test contains some simple utils for more readable tests.
// It is based on 'https://github.com/benbjohnson/testing'.
package test

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
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

type ExpectedErrorResult int

const (
	IsCorrectNil ExpectedErrorResult = iota
	IsCorrectErr
	IsWrongErr
)

// ExpectError checks if the actual error is expected given error, using errors.Is.
// If the expected error is nil, it returns IsCorrectNil if the actual error is also nil.
// Else it returns IsWrongErr.
// If the expected error is not nil, it returns IsCorrectErr if the actual error matches.
// Else it returns IsWrongErr.
//
// In all cases the test Fails if expected does not match actual.
//
// The three states are useful if you want to check for errors and do different things based on the result.
// For example if the actual error is only in some cases not nil and you want to continue execution in that case,
// but stop on the other cases.
func ExpectedError(tb testing.TB, exp, act error) ExpectedErrorResult {
	if exp == nil {
		if Ok(tb, act) {
			return IsCorrectNil
		}
		fmt.Printf("expected NO error, got \n%v", act)
		tb.Fail()
		return IsWrongErr
	}

	if errors.Is(act, exp) {
		return IsCorrectErr
	}

	fmt.Printf("expected error \n%v, got \n%v", exp, act)
	tb.Fail()
	return IsWrongErr
}
