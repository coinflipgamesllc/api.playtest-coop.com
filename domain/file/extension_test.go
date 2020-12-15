package file

import "testing"

func TestExtractExtension(t *testing.T) {
	var tests = []struct {
		filename          string
		expectedExtension Extension
	}{
		{"test.jpg", Extension("jpg")},
		{"test.JPEG", Extension("jpeg")},
		{"multiple.dots.in.name", Extension("name")},
		{".env", Extension("env")},
		{"no-extension", Extension("")},
	}

	for _, tt := range tests {
		actualExtension := ExtractExtension(tt.filename)
		if actualExtension != tt.expectedExtension {
			t.Errorf("got extension '%s', but expected '%s' for filename '%s'", actualExtension, tt.expectedExtension, tt.filename)
		}
	}
}

func TestExtensionsContain(t *testing.T) {
	var tests = []struct {
		category       Extensions
		extension      Extension
		expectContains bool
	}{
		{Images, Extension("jpg"), true},
		{Images, Extension("exe"), false},
		{Documents, Extension("pdf"), true},
		{Documents, Extension("docx"), false},
	}

	for _, tt := range tests {
		actual := tt.category.Contains(tt.extension)
		if actual != tt.expectContains {
			t.Errorf("extension '%s' should not be in '%v'", tt.extension, tt.category)
		}
	}
}
