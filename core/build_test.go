package core_test

import (
	"log"
	"testing"

	"github.com/spf13/afero"
	"github.com/verless/verless/core"
	"github.com/verless/verless/test"
)

const (
	outTestPath       = "../output-dir"
	projectFolderPath = "../example"
)

// TestRunFullBuild tests a full verless build and asserts
// that no errors arise.
func TestRunFullBuild(t *testing.T) {
	o := core.BuildOptions{
		OutputDir: outTestPath,
		Overwrite: true,
	}

	memMapFs := afero.NewMemMapFs()

	build, err := core.NewBuild(memMapFs, "../example", o)
	test.Ok(t, err)

	err = build.Run()
	test.Ok(t, err)

	if err := memMapFs.RemoveAll(outTestPath); err != nil {
		log.Println(err)
	}
}
