package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGenerateSlug kiểm tra generate slug từ title
func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  string
	}{
		{
			name:  "simple title",
			title: "Hello World",
			want:  "hello-world",
		},
		{
			name:  "title with special characters",
			title: "Hello World!",
			want:  "hello-world",
		},
		{
			name:  "title with multiple spaces",
			title: "Hello    World",
			want:  "hello-world",
		},
		{
			name:  "title with dashes",
			title: "Hello-World",
			want:  "hello-world",
		},
		{
			name:  "title with numbers",
			title: "Article 123",
			want:  "article-123",
		},
		{
			name:  "title with special chars only",
			title: "!!!",
			want:  "article", // Fallback to "article"
		},
		{
			name:  "empty title",
			title: "",
			want:  "article", // Fallback to "article"
		},
		{
			name:  "title with unicode",
			title: "Xin chào Việt Nam",
			want:  "xin-cho-vit-nam", // Unicode characters may be removed or transformed
		},
		{
			name:  "title with trailing spaces",
			title: "  Hello World  ",
			want:  "hello-world",
		},
		{
			name:  "title with uppercase",
			title: "HELLO WORLD",
			want:  "hello-world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateSlug(tt.title)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestGenerateUniqueSlug kiểm tra generate unique slug
func TestGenerateUniqueSlug(t *testing.T) {
	tests := []struct {
		name         string
		baseSlug     string
		isSlugExists func(string) bool
		want         string
	}{
		{
			name:     "slug not exists",
			baseSlug: "hello-world",
			isSlugExists: func(s string) bool {
				return false
			},
			want: "hello-world",
		},
		{
			name:     "slug exists, need to add number",
			baseSlug: "hello-world",
			isSlugExists: func(s string) bool {
				return s == "hello-world"
			},
			want: "hello-world-1",
		},
		{
			name:     "slug exists with multiple numbers",
			baseSlug: "hello-world",
			isSlugExists: func(s string) bool {
				// hello-world và hello-world-1 đều tồn tại
				return s == "hello-world" || s == "hello-world-1"
			},
			want: "hello-world-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateUniqueSlug(tt.baseSlug, tt.isSlugExists)
			assert.Equal(t, tt.want, got)
		})
	}
}

