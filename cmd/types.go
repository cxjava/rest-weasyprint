package main

const (
	NoShare FileShareService = ""        // No sharing, return PDF directly
	FileIO  FileShareService = "file.io" // https://file.io
	KITC    FileShareService = "ki.tc"   // https://ki.tc
	CVSH    FileShareService = "c-v.sh"  // https://c-v.sh
)

// HTMLRequest represents JSON request structure
type HTMLRequest struct {
	HTML         string                 `json:"html"`
	Options      map[string]interface{} `json:"options,omitempty"`
	ShareService string                 `json:"share_service,omitempty"`
}

// WeasyPrintOptions represents weasyprint supported options
type WeasyPrintOptions struct {
	// Basic options
	Encoding  string `json:"encoding"`
	MediaType string `json:"media_type"`
	BaseURL   string `json:"base_url"`

	// PDF related options
	PDFIdentifier string `json:"pdf_identifier"`
	PDFVariant    string `json:"pdf_variant"`
	PDFVersion    string `json:"pdf_version"`
	PDFForms      bool   `json:"pdf_forms"`

	// Output options
	UncompressedPDF     bool `json:"uncompressed_pdf"`
	CustomMetadata      bool `json:"custom_metadata"`
	PresentationalHints bool `json:"presentational_hints"`
	SRGB                bool `json:"srgb"`
	OptimizeImages      bool `json:"optimize_images"`
	FullFonts           bool `json:"full_fonts"`
	Hinting             bool `json:"hinting"`

	// Quality and performance options
	JPEGQuality int    `json:"jpeg_quality"`
	DPI         int    `json:"dpi"`
	CacheFolder string `json:"cache_folder"`
	Timeout     int    `json:"timeout"`

	// Log level
	Verbose bool `json:"verbose"`
	Debug   bool `json:"debug"`
	Quiet   bool `json:"quiet"`
}

// UploadedFileInfo stores uploaded file information
type UploadedFileInfo struct {
	HTMLPath     string
	CSSPaths     []string
	Attachments  []string
	Options      *WeasyPrintOptions
	Filename     string           // Add filename field
	ShareService FileShareService // Sharing service
}

// FileShareService defines supported third-party file sharing services
type FileShareService string

// ShareResponse represents the response from file sharing service
type ShareResponse struct {
	Link     string `json:"link"`
	Service  string `json:"service"`
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
	Filename string `json:"filename,omitempty"`
}
