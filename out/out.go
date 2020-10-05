// Package out provides functions for printing output to the user.
//
// Instead of using the standard library, this package should be used
// for printing user information and logging in general. Note that the
// out package isn't optimized for hot paths at the moment.
//
// This package is inspired by https://github.com/kubernetes/minikube.
package out

import (
	"fmt"
	"io"
	"os"

	"github.com/verless/verless/out/style"
)

var (
	outFile io.Writer = os.Stdout
	errFile io.Writer = os.Stderr
)

// T prints a prefixed, formatted text to the out file as a new line.
func T(prefix style.Emoji, format string, a ...interface{}) {
	printf(outFile, prefix, format, a...)
}

// Err prints a prefixed, formatted text to the error file as a new line.
func Err(prefix style.Emoji, format string, a ...interface{}) {
	printf(errFile, prefix, format, a...)
}

func printf(w io.Writer, prefix style.Emoji, format string, a ...interface{}) {
	formatted := fmt.Sprintf(format, a...)
	var output string

	if prefix == style.None {
		output = fmt.Sprintf("%s\n", formatted)
	} else {
		output = fmt.Sprintf("%s %s\n", prefix, formatted)
	}

	_, _ = w.Write([]byte(output))
}
