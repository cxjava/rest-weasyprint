package main

import (
	"log"
)

// PDFService encapsulates PDF generation related logic
type PDFService struct {
	logger *log.Logger
}

// NewPDFService creates a new PDF service instance
func NewPDFService(logger *log.Logger) *PDFService {
	return &PDFService{
		logger: logger,
	}
}
