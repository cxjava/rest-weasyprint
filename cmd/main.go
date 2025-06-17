package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	logger := log.New(os.Stdout, "PDF-SERVICE: ", log.LstdFlags|log.Lshortfile)

	router := setupRouter()
	pdfService := NewPDFService(logger)

	// Register routes
	registerRoutes(router, pdfService)

	// Get current timeout setting
	timeoutSeconds := getTimeoutFromEnv()

	logger.Printf("Server started, listening on port %s, request timeout %d seconds", DefaultPort, timeoutSeconds)
	if err := http.ListenAndServe(DefaultPort, router); err != nil {
		logger.Fatalf("Server startup failed: %v", err)
	}
}

// setupRouter configures router middleware
func setupRouter() *chi.Mux {
	router := chi.NewRouter()

	// Get timeout setting from environment variables
	timeoutSeconds := getTimeoutFromEnv()
	timeout := time.Duration(timeoutSeconds) * time.Second

	// Add middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(timeout))

	return router
}

// registerRoutes registers all routes
func registerRoutes(router *chi.Mux, service *PDFService) {
	// Health check
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok", "message": "PDF generation service is running"}`))
	})

	// API route group
	router.Route("/api/v1/pdf", func(r chi.Router) {
		r.Post("/render/file", service.HandleFileUpload) // File upload rendering
		r.Post("/render/html", service.HandleHTMLRender) // HTML string rendering
		r.Get("/render/html", service.HandleHTMLRender)  // GET test interface
		r.Get("/version", versionHandler)                // Version information
	})
}
