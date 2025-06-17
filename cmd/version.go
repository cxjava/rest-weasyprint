package main

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
	"sync"
)

var (
	Version            = "dev"
	Commit             = "none"
	RepoUrl            = "unknown"
	BuildDate          = "unknown"
	BuiltBy            = "unknown"
	BuiltWithGoVersion = "unknown"
)

// Only cache weasyprintVersion
var (
	weasyprintVersion string
	once              sync.Once
	initErr           error
)

func versionHandler(w http.ResponseWriter, r *http.Request) {
	// need to get weasyprint version dynamically
	once.Do(func() {
		cmd := exec.Command("weasyprint", "--version")
		output, err := cmd.Output()
		if err != nil {
			initErr = err
			return
		}
		weasyprintVersion = strings.TrimSpace(string(output))
	})
	if initErr != nil {
		http.Error(w, "failed to get weasyprint version: "+initErr.Error(), http.StatusInternalServerError)
		return
	}
	info := map[string]string{
		"apiVersion":         Version,
		"buildDate":          BuildDate,
		"builtBy":            BuiltBy,
		"builtWithGoVersion": BuiltWithGoVersion,
		"commit":             Commit,
		"repoUrl":            RepoUrl,
		"weasyprintVersion":  weasyprintVersion,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}
