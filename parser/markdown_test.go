package parser

import (
	"github.com/verless/verless/test"
	"testing"
	"time"
)

// TestMarkdown_ParsePage checks if a parsed Markdown file is
// converted to a model.Page instance correctly.
func TestMarkdown_ParsePage(t *testing.T) {
	parser := NewMarkdown()
	tests := []struct {
		src     string
		title   string
		date    time.Time
		tags    []string
		content string
	}{
		{
			src: `---
Title: Coffee Roasting Basics
Date: 2020-03-30
Tags:
    - Coffee
    - Roasting
---

This is a blog post.`,
			title:   "Coffee Roasting Basics",
			date:    time.Time{},
			tags:    []string{"Coffee", "Roasting"},
			content: "<p>This is a blog post.</p>\n",
		},
	}

	for _, testCase := range tests {
		page, err := parser.ParsePage([]byte(testCase.src))
		if err != nil {
			t.Fatal(err)
		}

		test.Equals(t, testCase.title, page.Title)
		test.Equals(t, testCase.tags, page.Tags)
		test.Equals(t, testCase.content, page.Content)
	}
}
