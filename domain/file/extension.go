package file

import "strings"

// Extension indicates the file type
type Extension string

// ExtractExtension will pull the extension from a filename
func ExtractExtension(filename string) Extension {
	idx := strings.LastIndex(filename, ".")
	if idx == -1 {
		return ""
	}

	return Extension(strings.ToLower(filename[idx+1:]))
}

// Extensions is a collection of extensions, probably of a related type (images)
type Extensions []Extension

// Contains checks if a given extension is in the list of extensions
func (e Extensions) Contains(ext Extension) bool {
	for _, x := range e {
		if x == ext {
			return true
		}
	}

	return false
}

var (
	// Images extensions we allow
	Images = Extensions{"png", "jpg", "jpeg", "svg", "tif", "tiff"}

	// Documents extensions we allow
	Documents = Extensions{"pdf"}
)
