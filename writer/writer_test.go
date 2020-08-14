// File writer_test.go tests the writer.

package writer

import (
	"github.com/verless/verless/util/test"
	"os"
	"path"
	"testing"
)

const (
	testPath    = "../example"
	testOutPath = "../test-out-path"
)

func TestWriter_removeOutDirIfPermitted(t *testing.T) {

	tests := map[string]struct {
		overwrite bool
		// beforeTest is a callback which creates the folders / files to test a specific testcase.
		beforeTest func()
		// cleanupTest is a callback which the folders / files created by beforeTest and the test itself.
		cleanupTest func()

		expectedError string
	}{
		"normal": {
			overwrite:   false,
			beforeTest:  func() {},
			cleanupTest: func() {},
		},
		"already exists": {
			overwrite: false,
			beforeTest: func() {
				test.Ok(t, os.Mkdir(testOutPath, os.ModePerm))

				_, err := os.Create(path.Join(testOutPath, "anyFile.txt"))
				test.Ok(t, err)
			},
			cleanupTest: func() {
				err := os.RemoveAll(testOutPath)
				test.Ok(t, err)
			},
			expectedError: "the output folder already exists and is not empty\ncondisder using the --overwrite flag",
		},
		"already exists but without file": {
			overwrite: false,
			beforeTest: func() {
				test.Ok(t, os.Mkdir(testOutPath, os.ModePerm))
			},
			cleanupTest: func() {
				err := os.RemoveAll(testOutPath)
				test.Ok(t, err)
			},
		},
		"already exists with overwrite": {
			overwrite: true,
			beforeTest: func() {
				test.Ok(t, os.Mkdir(testOutPath, os.ModePerm))
				_, err := os.Create(path.Join(testOutPath, "anyFile.txt"))
				test.Ok(t, err)
			},
			cleanupTest: func() {
				err := os.RemoveAll(testOutPath)
				test.Ok(t, err)
			},
		},
		"not exists with overwrite": {
			overwrite:   true,
			beforeTest:  func() {},
			cleanupTest: func() {},
		},
	}

	for caseName, testCase := range tests {
		t.Logf("Testing '%s'", caseName)

		w := setupNewWriter(t, testCase.overwrite)

		testCase.beforeTest()

		err := w.removeOutDirIfPermitted()

		if testCase.expectedError == "" {
			test.Ok(t, err)
		} else {
			test.Equals(t, testCase.expectedError, err.Error())
		}

		testCase.cleanupTest()
	}
}

func setupNewWriter(t testing.TB, overwrite bool) *writer {
	w, err := New(testPath, testOutPath, overwrite)
	if err != nil {
		t.Errorf("New should not throw an error: %v", err)
	}
	return w
}
