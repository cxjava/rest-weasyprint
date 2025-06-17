package main

import (
	"os"
	"strconv"
)

// Configuration constants
const (
	DefaultPort           = ":8080"
	DefaultTimeoutSeconds = 30       // Default timeout in seconds
	MaxUploadSize         = 32 << 20 // 32MB
	DefaultPageMargin     = "2cm 2.5cm"
	DefaultPageSize       = "A4"
)

// getTimeoutFromEnv gets timeout setting from environment variables, uses default if not exists
func getTimeoutFromEnv() int {
	timeoutSeconds := DefaultTimeoutSeconds
	if envTimeout := os.Getenv("WEB_TIME_OUT_SECOND"); envTimeout != "" {
		if parsedTimeout, err := strconv.Atoi(envTimeout); err == nil && parsedTimeout > 0 {
			timeoutSeconds = parsedTimeout
		}
	}
	return timeoutSeconds
}
