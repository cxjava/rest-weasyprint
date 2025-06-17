package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

// HandleFileUpload handles file upload and generates PDF
func (s *PDFService) HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	// Validate request type
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		http.Error(w, "multipart/form-data request required", http.StatusBadRequest)
		return
	}

	// Parse form
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		s.logger.Printf("Failed to parse form: %v", err)
		http.Error(w, "Form parsing failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Create temporary directory
	tempDir, err := s.createTempDir()
	if err != nil {
		s.logger.Printf("Failed to create temporary directory: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer s.cleanupTempDir(tempDir)

	// Process uploaded files
	fileInfo, err := s.processUploadedFiles(r, tempDir)
	if err != nil {
		s.logger.Printf("Failed to process uploaded files: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If sharing is needed, use temporary file instead of writing directly to response
	if fileInfo.ShareService != NoShare {
		tempFile, err := os.CreateTemp("", "pdfshare-*.pdf")
		if err != nil {
			s.logger.Printf("Failed to create temporary PDF file: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer os.Remove(tempFile.Name())

		// Generate PDF to temporary file
		if err := s.generatePDFFromFiles(r.Context(), tempFile, fileInfo); err != nil {
			s.logger.Printf("PDF generation failed: %v", err)
			http.Error(w, "PDF generation failed", http.StatusInternalServerError)
			return
		}
		tempFile.Close()

		// Upload to sharing service
		response, err := s.uploadToShareService(tempFile.Name(), fileInfo.Filename, fileInfo.ShareService)
		if err != nil {
			s.logger.Printf("Failed to upload to sharing service: %v", err)
			http.Error(w, "Failed to upload to sharing service: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Regular response, return PDF directly
	s.setPDFHeaders(w, fileInfo.Filename)

	// Generate PDF
	if err := s.generatePDFFromFiles(r.Context(), w, fileInfo); err != nil {
		s.logger.Printf("PDF generation failed: %v", err)
		http.Error(w, "PDF generation failed", http.StatusInternalServerError)
		return
	}
}

// HandleHTMLRender handles HTML string rendering
func (s *PDFService) HandleHTMLRender(w http.ResponseWriter, r *http.Request) {
	var htmlContent string
	var filename string
	var options *WeasyPrintOptions
	var shareService FileShareService = NoShare

	if r.Method == "POST" {
		// Handle JSON request
		var req HTMLRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON format error: "+err.Error(), http.StatusBadRequest)
			return
		}
		htmlContent = req.HTML

		// Get filename
		filename = r.URL.Query().Get("filename")
		if filename == "" {
			filename = "document.pdf" // Default value
		}

		// Process options
		if req.Options != nil {
			options = s.validateOptions(req.Options)
		} else {
			options = getDefaultOptions()
		}

		// Handle sharing service
		shareServiceParam := req.ShareService
		if shareServiceParam == "" {
			shareServiceParam = r.URL.Query().Get("share_service")
		}

		switch shareServiceParam {
		case string(FileIO):
			shareService = FileIO
		case string(KITC):
			shareService = KITC
		case string(CVSH):
			shareService = CVSH
		}

	} else {
		// GET request, return example
		htmlContent = `
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Test Document</title>
			<style>
				body { font-family: Arial, sans-serif; margin: 20px; }
				h1 { color: #333; }
			</style>
		</head>
		<body>
			<h1>Hello, World! 🌍</h1>
			<p>This is a test PDF document.</p>
		</body>
		</html>`
		filename = r.URL.Query().Get("filename")
		if filename == "" {
			filename = "test.pdf" // Default value
		}
		options = getDefaultOptions()

		// Handle sharing service
		shareServiceParam := r.URL.Query().Get("share_service")
		switch shareServiceParam {
		case string(FileIO):
			shareService = FileIO
		case string(KITC):
			shareService = KITC
		case string(CVSH):
			shareService = CVSH
		}
	}

	// If sharing is needed, use temporary file instead of writing directly to response
	if shareService != NoShare {
		tempFile, err := os.CreateTemp("", "pdfshare-*.pdf")
		if err != nil {
			s.logger.Printf("Failed to create temporary PDF file: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer os.Remove(tempFile.Name())

		// Generate PDF to temporary file
		if err := s.generatePDFFromHTML(r.Context(), tempFile, htmlContent, options); err != nil {
			s.logger.Printf("PDF generation failed: %v", err)
			http.Error(w, "PDF generation failed", http.StatusInternalServerError)
			return
		}
		tempFile.Close()

		// Upload to sharing service
		response, err := s.uploadToShareService(tempFile.Name(), filename, shareService)
		if err != nil {
			s.logger.Printf("Failed to upload to sharing service: %v", err)
			http.Error(w, "Failed to upload to sharing service: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Regular response, return PDF directly
	s.setPDFHeaders(w, filename)

	// Generate PDF
	if err := s.generatePDFFromHTML(r.Context(), w, htmlContent, options); err != nil {
		s.logger.Printf("PDF generation failed: %v", err)
		http.Error(w, "PDF generation failed", http.StatusInternalServerError)
		return
	}
}
