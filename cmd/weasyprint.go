package main

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// buildWeasyPrintArgs builds weasyprint command arguments
func (s *PDFService) buildWeasyPrintArgs(options *WeasyPrintOptions) []string {
	var args []string

	if options == nil {
		return args
	}

	// Basic options
	if options.Encoding != "" {
		args = append(args, "--encoding", options.Encoding)
	}

	if options.MediaType != "" && options.MediaType != "print" {
		args = append(args, "--media-type", options.MediaType)
	}

	if options.BaseURL != "" {
		args = append(args, "--base-url", options.BaseURL)
	}

	// PDF related options
	if options.PDFIdentifier != "" {
		args = append(args, "--pdf-identifier", options.PDFIdentifier)
	}

	if options.PDFVariant != "" {
		args = append(args, "--pdf-variant", options.PDFVariant)
	}

	if options.PDFVersion != "" {
		args = append(args, "--pdf-version", options.PDFVersion)
	}

	if options.PDFForms {
		args = append(args, "--pdf-forms")
	}

	// Output options
	if options.UncompressedPDF {
		args = append(args, "--uncompressed-pdf")
	}

	if options.CustomMetadata {
		args = append(args, "--custom-metadata")
	}

	if options.PresentationalHints {
		args = append(args, "--presentational-hints")
	}

	if options.SRGB {
		args = append(args, "--srgb")
	}

	if options.OptimizeImages {
		args = append(args, "--optimize-images")
	}

	if options.FullFonts {
		args = append(args, "--full-fonts")
	}

	if options.Hinting {
		args = append(args, "--hinting")
	}

	// Quality and performance options
	if options.JPEGQuality > 0 && options.JPEGQuality != 80 {
		args = append(args, "--jpeg-quality", strconv.Itoa(options.JPEGQuality))
	}

	if options.DPI > 0 && options.DPI != 96 {
		args = append(args, "--dpi", strconv.Itoa(options.DPI))
	}

	if options.CacheFolder != "" {
		args = append(args, "--cache-folder", options.CacheFolder)
	}

	if options.Timeout > 0 && options.Timeout != 30 {
		args = append(args, "--timeout", strconv.Itoa(options.Timeout))
	}

	// Log level
	if options.Verbose {
		args = append(args, "--verbose")
	}

	if options.Debug {
		args = append(args, "--debug")
	}

	if options.Quiet {
		args = append(args, "--quiet")
	}

	return args
}

// executeWeasyPrint executes weasyprint command
func (s *PDFService) executeWeasyPrint(ctx context.Context, w io.Writer, args []string) error {
	s.logger.Printf("Executing weasyprint command: %v", args)

	cmd := exec.CommandContext(ctx, "weasyprint", args...)
	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("weasyprint execution failed: %v", err)
	}

	return nil
}

// generatePDFFromFiles generates PDF from files
func (s *PDFService) generatePDFFromFiles(ctx context.Context, w io.Writer, fileInfo *UploadedFileInfo) error {
	// Build weasyprint command arguments
	args := s.buildWeasyPrintArgs(fileInfo.Options)

	// Add CSS files
	for _, cssPath := range fileInfo.CSSPaths {
		args = append(args, "--stylesheet", cssPath)
	}

	// Add attachments
	for _, attachment := range fileInfo.Attachments {
		args = append(args, "--attachment", attachment)
	}

	// Add HTML file and output
	args = append(args, fileInfo.HTMLPath, "-")

	return s.executeWeasyPrint(ctx, w, args)
}

// generatePDFFromHTML generates PDF from HTML string
func (s *PDFService) generatePDFFromHTML(ctx context.Context, w io.Writer, htmlContent string, options *WeasyPrintOptions) error {
	// Check if htmlContent is a URL
	isURL := false
	if strings.HasPrefix(htmlContent, "http://") || strings.HasPrefix(htmlContent, "https://") {
		// Try to parse URL to ensure format is correct
		_, err := url.Parse(htmlContent)
		if err == nil {
			isURL = true
		}
	}

	// Build command arguments
	args := s.buildWeasyPrintArgs(options)

	if isURL {
		// If it's a URL, add directly to arguments
		s.logger.Printf("Detected URL: %s", htmlContent)
		args = append(args, htmlContent, "-")
	} else {
		// Not a URL, create temporary HTML file
		tempFile, err := os.CreateTemp("", "*.html")
		if err != nil {
			return fmt.Errorf("failed to create temporary file: %v", err)
		}

		defer func() {
			tempFile.Close()
			os.Remove(tempFile.Name())
		}()

		// Write HTML content
		if _, err := tempFile.WriteString(htmlContent); err != nil {
			return fmt.Errorf("failed to write HTML content: %v", err)
		}

		if err := tempFile.Close(); err != nil {
			return fmt.Errorf("failed to close temporary file: %v", err)
		}

		args = append(args, tempFile.Name(), "-")
	}

	return s.executeWeasyPrint(ctx, w, args)
}

// getDefaultOptions returns default weasyprint options
func getDefaultOptions() *WeasyPrintOptions {
	return &WeasyPrintOptions{
		MediaType: "print",
	}
}

// validateOptions validates and cleans options
func (s *PDFService) validateOptions(options map[string]interface{}) *WeasyPrintOptions {
	result := getDefaultOptions()

	// Define safe options list (exclude unsafe or unsupported options)
	safeOptions := map[string]bool{
		"encoding": true, "media_type": true, "base_url": true,
		"pdf_identifier": true, "pdf_variant": true, "pdf_version": true, "pdf_forms": true,
		"uncompressed_pdf": true, "custom_metadata": true, "presentational_hints": true,
		"srgb": true, "optimize_images": true, "full_fonts": true, "hinting": true,
		"jpeg_quality": true, "dpi": true, "timeout": true,
		"verbose": true, "debug": true, "quiet": true,
	}

	for key, value := range options {
		if !safeOptions[key] {
			s.logger.Printf("Ignoring unsafe option: %s", key)
			continue
		}

		switch key {
		case "encoding":
			if str, ok := value.(string); ok {
				result.Encoding = str
			}
		case "media_type":
			if str, ok := value.(string); ok {
				result.MediaType = str
			}
		case "base_url":
			if str, ok := value.(string); ok {
				result.BaseURL = str
			}
		case "pdf_identifier":
			if str, ok := value.(string); ok {
				result.PDFIdentifier = str
			}
		case "pdf_variant":
			if str, ok := value.(string); ok && s.isValidPDFVariant(str) {
				result.PDFVariant = str
			}
		case "pdf_version":
			if str, ok := value.(string); ok {
				result.PDFVersion = str
			}
		case "pdf_forms":
			if b, ok := value.(bool); ok {
				result.PDFForms = b
			}
		case "uncompressed_pdf":
			if b, ok := value.(bool); ok {
				result.UncompressedPDF = b
			}
		case "custom_metadata":
			if b, ok := value.(bool); ok {
				result.CustomMetadata = b
			}
		case "presentational_hints":
			if b, ok := value.(bool); ok {
				result.PresentationalHints = b
			}
		case "srgb":
			if b, ok := value.(bool); ok {
				result.SRGB = b
			}
		case "optimize_images":
			if b, ok := value.(bool); ok {
				result.OptimizeImages = b
			}
		case "full_fonts":
			if b, ok := value.(bool); ok {
				result.FullFonts = b
			}
		case "hinting":
			if b, ok := value.(bool); ok {
				result.Hinting = b
			}
		case "jpeg_quality":
			if quality := s.parseIntValue(value, 0, 95); quality > 0 {
				result.JPEGQuality = quality
			}
		case "dpi":
			if dpi := s.parseIntValue(value, 50, 600); dpi > 0 {
				result.DPI = dpi
			}
		case "timeout":
			if timeout := s.parseIntValue(value, 1, 300); timeout > 0 {
				result.Timeout = timeout
			}
		case "verbose":
			if b, ok := value.(bool); ok {
				result.Verbose = b
			}
		case "debug":
			if b, ok := value.(bool); ok {
				result.Debug = b
			}
		case "quiet":
			if b, ok := value.(bool); ok {
				result.Quiet = b
			}
		}
	}

	return result
}

// parseIntValue parses integer value and validates range
func (s *PDFService) parseIntValue(value interface{}, min, max int) int {
	var intVal int

	switch v := value.(type) {
	case int:
		intVal = v
	case float64:
		intVal = int(v)
	case string:
		if parsed, err := strconv.Atoi(v); err == nil {
			intVal = parsed
		} else {
			return 0
		}
	default:
		return 0
	}

	if intVal < min || intVal > max {
		s.logger.Printf("Integer value %d out of range [%d, %d]", intVal, min, max)
		return 0
	}

	return intVal
}

// isValidPDFVariant validates if PDF variant is valid
func (s *PDFService) isValidPDFVariant(variant string) bool {
	validVariants := []string{
		"pdf/a-1b", "pdf/a-2b", "pdf/a-3b", "pdf/a-4b",
		"pdf/a-2u", "pdf/a-3u", "pdf/a-4u",
		"pdf/ua-1", "debug",
	}

	for _, valid := range validVariants {
		if variant == valid {
			return true
		}
	}

	s.logger.Printf("Invalid PDF variant: %s", variant)
	return false
}
