package domain

import (
	"math/rand"
	"testing"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain/file"
)

func TestNewImage(t *testing.T) {
	var tests = []struct {
		extension     string
		expectedError error
	}{
		{"png", nil},
		{"exe", file.InvalidExtension{}},
	}

	for _, tt := range tests {
		f, err := NewImage(User{ID: 123}, "file."+tt.extension, "bucket", "object", 123)
		if tt.expectedError != nil && err == nil {
			t.Errorf("Dissallowed extension allowed for new image: %s", tt.extension)
		}

		if tt.expectedError == nil && err != nil {
			t.Errorf("Allowed extension not allowed for new image: %s", tt.extension)
		}

		if tt.expectedError != nil && err != nil {
			if _, ok := err.(file.InvalidExtension); !ok {
				t.Error("Error should be InvalidExtension")
			}
		} else {
			if f.UploadedBy.ID != 123 {
				t.Error("UploadedBy not set for new image")
			}

			if f.Filename != "file."+tt.extension {
				t.Error("Filename not set for new image")
			}

			if f.Bucket != "bucket" {
				t.Error("Bucket not set for new image")
			}

			if f.Object != "object" {
				t.Error("Object not set for new image")
			}

			if f.Size != 123 {
				t.Error("Size not set for new image")
			}

			if f.URL != "https://assets.playtest-coop.com/object" {
				t.Error("URL decorator not run on new image")
			}
		}
	}
}

func TestNewSellSheet(t *testing.T) {
	var tests = []struct {
		extension     string
		expectedError error
	}{
		{"pdf", nil},
		{"exe", file.InvalidExtension{}},
	}

	for _, tt := range tests {
		f, err := NewSellSheet(User{ID: 123}, "file."+tt.extension, "bucket", "object", 123)
		if tt.expectedError != nil && err == nil {
			t.Errorf("Dissallowed extension allowed for new sellsheet: %s", tt.extension)
		}

		if tt.expectedError == nil && err != nil {
			t.Errorf("Allowed extension not allowed for new sellsheet: %s", tt.extension)
		}

		if tt.expectedError != nil && err != nil {
			if _, ok := err.(file.InvalidExtension); !ok {
				t.Error("Error should be InvalidExtension")
			}
		} else {
			if f.UploadedBy.ID != 123 {
				t.Error("UploadedBy not set for new sellsheet")
			}

			if f.Filename != "file."+tt.extension {
				t.Error("Filename not set for new sellsheet")
			}

			if f.Bucket != "bucket" {
				t.Error("Bucket not set for new sellsheet")
			}

			if f.Object != "object" {
				t.Error("Object not set for new sellsheet")
			}

			if f.Size != 123 {
				t.Error("Size not set for new sellsheet")
			}

			if f.URL != "https://assets.playtest-coop.com/object" {
				t.Error("URL decorator not run on new sellsheet")
			}
		}
	}
}

func TestNewPrintAndPlay(t *testing.T) {
	var tests = []struct {
		extension     string
		expectedError error
	}{
		{"pdf", nil},
		{"exe", file.InvalidExtension{}},
	}

	for _, tt := range tests {
		f, err := NewPrintAndPlay(User{ID: 123}, "file."+tt.extension, "bucket", "object", 123)
		if tt.expectedError != nil && err == nil {
			t.Errorf("Dissallowed extension allowed for new pnp: %s", tt.extension)
		}

		if tt.expectedError == nil && err != nil {
			t.Errorf("Allowed extension not allowed for new pnp: %s", tt.extension)
		}

		if tt.expectedError != nil && err != nil {
			if _, ok := err.(file.InvalidExtension); !ok {
				t.Error("Error should be InvalidExtension")
			}
		} else {
			if f.UploadedBy.ID != 123 {
				t.Error("UploadedBy not set for new pnp")
			}

			if f.Filename != "file."+tt.extension {
				t.Error("Filename not set for new pnp")
			}

			if f.Bucket != "bucket" {
				t.Error("Bucket not set for new pnp")
			}

			if f.Object != "object" {
				t.Error("Object not set for new pnp")
			}

			if f.Size != 123 {
				t.Error("Size not set for new pnp")
			}

			if f.URL != "https://assets.playtest-coop.com/object" {
				t.Error("URL decorator not run on new pnp")
			}
		}
	}
}

func TestUpdateCaption(t *testing.T) {
	var tests = []struct {
		file            *File
		newCaption      string
		expectedCaption string
	}{
		{&File{Caption: "Original caption"}, "New caption", "New caption"},
		{&File{Caption: "Original caption"}, "", "Original caption"},
		{&File{Caption: "Original caption"}, "Original caption", "Original caption"},
	}

	for _, tt := range tests {
		tt.file.UpdateCaption(tt.newCaption)
		actual := tt.file.Caption
		if tt.expectedCaption != actual {
			t.Errorf("UpdateCaption incorrect")
		}
	}
}

func TestBelongsTo(t *testing.T) {
	f := &File{}
	f.BelongsTo(&Game{ID: 123})

	if *f.GameID != uint(123) {
		t.Error("Belongs to needs to create the relationship with game")
	}
}

func TestGenerateObjectName(t *testing.T) {
	rand.Seed(42)

	name := "test-name.png"
	expected := "nux1a0S2Ts65XHQkWNH_WfLl8wzB8STvqutDHMwY-4Q.png"
	actual := GenerateObjectName(name, "png")

	if actual != expected {
		t.Errorf("Generated object name is incorrect. Expected '%s', got '%s'", actual, expected)
	}
}
