package core_test

import (
	"log"
	"os"
	"testing"

	"github.com/verless/verless/config"
	"github.com/verless/verless/core"
	"github.com/verless/verless/test"
)

const (
	outTestPath       = "../test-output-dir"
	projectFolderPath = "../example"
)

// TestRunFullBuild tests a full verless build and asserts
// that no errors arise.
func TestRunFullBuild(t *testing.T) {
	o := core.BuildOptions{
		OutputDir: outTestPath,
		Overwrite: true,
	}

	cfg, err := config.FromFile(projectFolderPath, config.Filename)
	test.Ok(t, err)

	errs := core.RunBuild("../example", o, cfg)
	for _, err := range errs {
		test.Ok(t, err)
	}

	if err := os.RemoveAll(outTestPath); err != nil {
		log.Println(err)
	}
}
