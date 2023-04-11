package tpl

import (
	"path/filepath"
	"testing"
	"text/template"

	"github.com/verless/verless/test"
	"github.com/verless/verless/theme"
)

const (
	projectPath = "../example"
	testKey     = "test key"
	invalidKey  = "invalid key"
)

// TestRegister checks if the Register function register all templates
// correctly and if it doesn't allow overwriting a template by default.
func TestRegister(t *testing.T) {
	// Don't use a map as the execution order is important here.
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
		pageTplPath := filepath.Join(theme.TemplatePath(projectPath, theme.Default), theme.PageTemplate)
		tplBasePath := filepath.Join(theme.TemplatePath(projectPath, theme.Default), theme.TemplateBase)

		_, err := Register(testCase.key, pageTplPath, tplBasePath, testCase.force)
		test.ExpectedError(t, testCase.expectedError, err)
	}
}

// TestGet checks if the Get function correctly returns a registered
// template.
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
