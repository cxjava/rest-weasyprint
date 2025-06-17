package main

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// processUploadedFiles processes uploaded files
func (s *PDFService) processUploadedFiles(r *http.Request, tempDir string) (*UploadedFileInfo, error) {
	form := r.MultipartForm
	fileInfo := &UploadedFileInfo{
		Options:      getDefaultOptions(),
		Filename:     "document.pdf", // Default filename
		ShareService: NoShare,        // Default no sharing
	}

	// Get custom filename
	if filenameValue := r.URL.Query().Get("filename"); len(filenameValue) > 0 {
		customFilename := strings.TrimSpace(filenameValue)
		if customFilename != "" {
			fileInfo.Filename = customFilename
			// Ensure filename ends with .pdf
			if !strings.HasSuffix(strings.ToLower(fileInfo.Filename), ".pdf") {
				fileInfo.Filename += ".pdf"
			}
		}
	}

	// Get sharing service parameter
	if shareServiceValue := r.URL.Query().Get("share_service"); len(shareServiceValue) > 0 {
		shareService := strings.TrimSpace(shareServiceValue)
		switch shareService {
		case string(FileIO):
			fileInfo.ShareService = FileIO
		case string(KITC):
			fileInfo.ShareService = KITC
		case string(CVSH):
			fileInfo.ShareService = CVSH
		}
	}

	// Iterate through all file fields
	for fieldName, files := range form.File {
		for _, fileHeader := range files {
			filePath, err := s.saveUploadedFile(fileHeader, tempDir)
			if err != nil {
				return nil, fmt.Errorf("failed to save file: %v", err)
			}

			// Classify files by field name
			switch {
			case fieldName == "html":
				fileInfo.HTMLPath = filePath
			case strings.HasPrefix(fieldName, "css."):
				fileInfo.CSSPaths = append(fileInfo.CSSPaths, filePath)
			case strings.HasPrefix(fieldName, "asset."):
				fileInfo.Attachments = append(fileInfo.Attachments, filePath)
			}
		}
	}

	// Process options field
	if optionsValues, exists := form.Value["options"]; exists && len(optionsValues) > 0 {
		var optionsMap map[string]interface{}
		if err := json.Unmarshal([]byte(optionsValues[0]), &optionsMap); err != nil {
			s.logger.Printf("Failed to parse options: %v", err)
			return nil, fmt.Errorf("invalid JSON format options")
		}
		fileInfo.Options = s.validateOptions(optionsMap)
	}

	// Validate required files
	if fileInfo.HTMLPath == "" {
		return nil, fmt.Errorf("missing HTML file")
	}

	// If no CSS files, create default one
	if len(fileInfo.CSSPaths) == 0 {
		defaultCSSPath, err := s.createDefaultCSS(tempDir)
		if err != nil {
			return nil, fmt.Errorf("failed to create default CSS: %v", err)
		}
		fileInfo.CSSPaths = append(fileInfo.CSSPaths, defaultCSSPath)
	}

	return fileInfo, nil
}

// saveUploadedFile saves uploaded file
func (s *PDFService) saveUploadedFile(fileHeader *multipart.FileHeader, tempDir string) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	dstPath := filepath.Join(tempDir, fileHeader.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to create target file: %v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file content: %v", err)
	}

	return dstPath, nil
}

// createTempDir creates temporary directory
func (s *PDFService) createTempDir() (string, error) {
	tempDir, err := os.MkdirTemp("", "pdfgen-*")
	if err != nil {
		return "", err
	}
	s.logger.Printf("Created temporary directory: %s", tempDir)
	return tempDir, nil
}

// cleanupTempDir cleans up temporary directory
func (s *PDFService) cleanupTempDir(tempDir string) {
	if err := os.RemoveAll(tempDir); err != nil {
		s.logger.Printf("Failed to cleanup temporary directory: %v", err)
	} else {
		s.logger.Printf("Cleaned up temporary directory: %s", tempDir)
	}
}

// createDefaultCSS creates default CSS file
func (s *PDFService) createDefaultCSS(tempDir string) (string, error) {
	defaultCSS := fmt.Sprintf("@page { size: %s; margin: %s; }", DefaultPageSize, DefaultPageMargin)
	cssPath := filepath.Join(tempDir, "default.css")

	if err := os.WriteFile(cssPath, []byte(defaultCSS), 0644); err != nil {
		return "", err
	}

	return cssPath, nil
}

// setPDFHeaders sets PDF response headers
func (s *PDFService) setPDFHeaders(w http.ResponseWriter, filename string) {
	w.Header().Set("Content-Type", "application/pdf")

	// Use url.PathEscape instead of url.QueryEscape
	// PathEscape encodes spaces as %20, not +
	encodedFilename := url.PathEscape(filename)

	// Set Content-Disposition header
	// filename= parameter for old browsers that don't support RFC 5987
	// filename*= parameter for modern browsers that support UTF-8
	disposition := fmt.Sprintf("attachment; filename=\"%s\"; filename*=UTF-8''%s",
		filename, encodedFilename)
	w.Header().Set("Content-Disposition", disposition)
}
