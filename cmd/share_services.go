package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

// uploadToShareService uploads PDF to third-party sharing service
func (s *PDFService) uploadToShareService(filePath, filename string, service FileShareService) (*ShareResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open PDF file: %v", err)
	}
	defer file.Close()

	switch service {
	case FileIO:
		return s.uploadToFileIO(file, filename)
	case KITC:
		return s.uploadToKITC(file, filename)
	case CVSH:
		return s.uploadToCVSH(file, filename)
	default:
		return nil, fmt.Errorf("unsupported sharing service: %s", service)
	}
}

// uploadToFileIO uploads to file.io service
func (s *PDFService) uploadToFileIO(file *os.File, filename string) (*ShareResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file content: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %v", err)
	}

	req, err := http.NewRequest("POST", "https://file.io", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	var fileIOResp struct {
		Success bool   `json:"success"`
		Key     string `json:"key"`
		Link    string `json:"link"`
		Expiry  string `json:"expiry"`
	}

	fmt.Println(resp.Status)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	fmt.Println(string(bodyBytes))

	err = json.Unmarshal(bodyBytes, &fileIOResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if !fileIOResp.Success {
		return nil, fmt.Errorf("failed to upload to file.io")
	}

	return &ShareResponse{
		Success:  true,
		Link:     fileIOResp.Link,
		Service:  string(FileIO),
		Filename: filename,
	}, nil
}

// uploadToKITC uploads to ki.tc service
func (s *PDFService) uploadToKITC(file *os.File, filename string) (*ShareResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file content: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %v", err)
	}

	req, err := http.NewRequest("POST", "https://ki.tc/file/u/", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to upload to ki.tc: %s", resp.Status)
	}
	// Parse JSON response
	var response struct {
		File struct {
			DownloadPage string `json:"download_page"`
		} `json:"file"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	if response.File.DownloadPage == "" {
		return nil, fmt.Errorf("missing download page URL in response")
	}

	return &ShareResponse{
		Success:  true,
		Link:     response.File.DownloadPage,
		Service:  string(KITC),
		Filename: filename,
	}, nil
}

// uploadToCVSH uploads to c-v.sh service
func (s *PDFService) uploadToCVSH(file *os.File, filename string) (*ShareResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file content: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %v", err)
	}

	req, err := http.NewRequest("POST", "https://c-v.sh", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to upload to c-v.sh: %s", resp.Status)
	}

	// c-v.sh returns URL directly as response body
	urlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	url := strings.TrimSpace(string(urlBytes))
	if !strings.HasPrefix(url, "http") {
		return nil, fmt.Errorf("invalid URL response: %s", url)
	}

	return &ShareResponse{
		Success:  true,
		Link:     url,
		Service:  string(CVSH),
		Filename: filename,
	}, nil
}
