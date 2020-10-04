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
)

// Unicode mappings for all supported emojis.
const (
	Tada = "\U0001f389"
)

var outFile io.Writer = os.Stdout

// T prints a prefixed, formatted text to outFile.
func T(prefix, format string, a ...interface{}) {
	output := fmt.Sprintf("%s %s\n", prefix, fmt.Sprintf(format, a))
	_, _ = outFile.Write([]byte(output))
}

// SetOutFile sets the output file, which defaults to os.Stdout.
func SetOutFile(w io.Writer) {
	outFile = w
}
