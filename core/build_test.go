package core_test

import (
	"os"
	"testing"

	"github.com/verless/verless/config"
	"github.com/verless/verless/core"
	"github.com/verless/verless/test"
)

const (
	// Yes, it is intention to use spaces just to test this also.
	outTestPath       = "out test path"
	projectFolderPath = "../example"
)

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

	_ = os.RemoveAll(outTestPath)
}
