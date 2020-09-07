package tpl

import (
	"path/filepath"
	"testing"
	"text/template"

	"github.com/verless/verless/config"
	"github.com/verless/verless/test"
)

const (
	projectFolderPath = "../example"
	testKey           = "test key"
	invalidKey        = "invalid key"
)

func TestRegister(t *testing.T) {
	// Don't use a map here as the execution order is important here.
	tests := []struct {
		testName      string
		force         bool
		expectedError error
		key           string
	}{
		{
			testName: "first template",
			key:      "test key",
		},
		{
			testName: "second template",
			key:      "test key2",
		},
		{
			testName:      "same template again",
			key:           "test key2",
			expectedError: ErrAlreadyRegistered,
		},
		{
			testName: "same template again with force true",
			key:      "test key2",
			force:    true,
		},
	}

	for _, testCase := range tests {
		t.Logf("Testing '%s'", testCase.testName)
		pageTplPath := filepath.Join(projectFolderPath, config.TemplateDir, config.PageTpl)

		_, err := Register(testCase.key, pageTplPath, testCase.force)
		test.ExpectedError(t, testCase.expectedError, err)
	}
}

func TestGet(t *testing.T) {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templates[testKey] = &template.Template{}

	tpl, err := Get(testKey)
	test.Ok(t, err)

	test.Assert(t, tpl == templates[testKey], "template has to be in map")

	_, err = Get(invalidKey)
	test.Assert(t, err != nil, "template key is invalid")
}
