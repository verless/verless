package parser

import (
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

	for _, test := range tests {
		page, err := parser.ParsePage([]byte(test.src))
		if err != nil {
			t.Fatal(err)
		}

		if page.Title != test.title {
			t.Errorf("expected title %s, got %s", test.title, page.Title)
		}

		if len(page.Tags) != len(test.tags) {
			t.Fatalf("expected %d tags, got %d", len(test.tags), len(page.Tags))
		}

		for i, tag := range page.Tags {
			if tag != test.tags[i] {
				t.Errorf("expected tag %s, got %s", test.tags[i], tag)
			}
		}

		if page.Content != test.content {
			t.Errorf("expected content %s, got %s", test.content, page.Content)
		}
	}
}
