package file

// Role tracks the purpose of a file
type Role string

const (
	// Image files (png, jpg, svg)
	Image Role = "Image"

	// SellSheet files (pdf)
	SellSheet = "SellSheet"

	// PrintAndPlay files (pdf)
	PrintAndPlay = "PrintAndPlay"
)
