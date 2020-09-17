package writer

import (
	"os"
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/verless/verless/fs"
	"github.com/verless/verless/test"
)

const (
	testPath    = "../example"
	testOutPath = "../test-out-path"
)

// TestWriter_removeOutDirIfExists check if the removeOutDirIfExists
// correctly removes a directory.
func TestWriter_removeOutDirIfExists(t *testing.T) {
	memMapFs := afero.NewMemMapFs()

	tests := map[string]struct {
		beforeTest    func()
		cleanupTest   func()
		expectedError string
	}{
		"normal": {
			beforeTest:  func() {},
			cleanupTest: func() {},
		},
		"already exists": {
			beforeTest: func() {
				test.Ok(t, memMapFs.Mkdir(testOutPath, os.ModePerm))

				file, err := memMapFs.Create(path.Join(testOutPath, "anyFile.txt"))
				test.Ok(t, err)
				_ = file.Close()
			},
			cleanupTest: func() {
				err := memMapFs.RemoveAll(testOutPath)
				test.Ok(t, err)
			},
		},
		"already exists but without file": {
			beforeTest: func() {
				test.Ok(t, memMapFs.Mkdir(testOutPath, os.ModePerm))
			},
			cleanupTest: func() {
				err := memMapFs.RemoveAll(testOutPath)
				test.Ok(t, err)
			},
		},
	}

	for caseName, testCase := range tests {
		t.Logf("Testing '%s'", caseName)

		w := setupNewWriter(memMapFs)

		testCase.beforeTest()

		err := fs.Rmdir(memMapFs, w.outputDir)

		if testCase.expectedError == "" {
			test.Ok(t, err)
		} else {
			test.Assert(t, err != nil && testCase.expectedError == err.Error(), "should error")
		}

		testCase.cleanupTest()
	}
}

// setupNewWriter initializes a new writer instance.
func setupNewWriter(fs afero.Fs) *writer {
	return New(fs, testPath, testOutPath, false)
}
