package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	baseURL    = "https://api.cloud.llamaindex.ai/api/v1/parsing"
	pollDelay  = 5 * time.Second
	maxRetries = 120 // 10 minutes with 5-second delays
)

type UploadResponse struct {
	ID string `json:"id"`
}

type JobStatus struct {
	Status string `json:"status"`
}

func main() {
	apiKey := os.Getenv("LLAMA_CLOUD_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: LLAMA_CLOUD_API_KEY environment variable not set")
		os.Exit(1)
	}

	// Create directories if they don't exist
	if err := os.MkdirAll("pdf", 0755); err != nil {
		fmt.Printf("Error creating pdf directory: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll("markdown", 0755); err != nil {
		fmt.Printf("Error creating markdown directory: %v\n", err)
		os.Exit(1)
	}

	// Get list of already processed files
	processed := getProcessedFiles()
	fmt.Printf("Found %d already processed files\n", len(processed))

	// Get list of PDF files to process
	pdfFiles, err := filepath.Glob("pdf/*.pdf")
	if err != nil {
		fmt.Printf("Error listing PDF files: %v\n", err)
		os.Exit(1)
	}

	if len(pdfFiles) == 0 {
		fmt.Println("No PDF files found in pdf/ directory")
		return
	}

	fmt.Printf("Found %d PDF files\n", len(pdfFiles))

	// Process each PDF
	for i, pdfPath := range pdfFiles {
		baseName := filepath.Base(pdfPath)
		baseNameNoExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))

		// Skip if already processed
		if processed[baseNameNoExt] {
			fmt.Printf("[%d/%d] Skipping %s (already processed)\n", i+1, len(pdfFiles), baseName)
			continue
		}

		fmt.Printf("[%d/%d] Processing %s...\n", i+1, len(pdfFiles), baseName)

		if err := processPDF(pdfPath, baseNameNoExt, apiKey); err != nil {
			fmt.Printf("Error processing %s: %v\n", baseName, err)
			continue
		}

		fmt.Printf("[%d/%d] ✓ Successfully processed %s\n", i+1, len(pdfFiles), baseName)
	}

	fmt.Println("\nAll files processed!")
}

func getProcessedFiles() map[string]bool {
	processed := make(map[string]bool)

	mdFiles, err := filepath.Glob("markdown/*.md")
	if err != nil {
		return processed
	}

	for _, mdPath := range mdFiles {
		baseName := filepath.Base(mdPath)
		baseNameNoExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))
		processed[baseNameNoExt] = true
	}

	return processed
}

func processPDF(pdfPath, baseName, apiKey string) error {
	// Step 1: Upload PDF
	jobID, err := uploadPDF(pdfPath, apiKey)
	if err != nil {
		return fmt.Errorf("upload failed: %w", err)
	}
	fmt.Printf("  Uploaded, job ID: %s\n", jobID)

	// Step 2: Wait for job to complete
	if err := waitForJob(jobID, apiKey); err != nil {
		return fmt.Errorf("job failed: %w", err)
	}
	fmt.Printf("  Job completed\n")

	// Step 3: Download markdown result
	markdown, err := downloadMarkdown(jobID, apiKey)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	// Step 4: Save to file
	outputPath := filepath.Join("markdown", baseName+".md")
	if err := os.WriteFile(outputPath, []byte(markdown), 0644); err != nil {
		return fmt.Errorf("save failed: %w", err)
	}

	return nil
}

func uploadPDF(pdfPath, apiKey string) (string, error) {
	file, err := os.Open(pdfPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(pdfPath))
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", baseURL+"/upload", body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var uploadResp UploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return "", err
	}

	return uploadResp.ID, nil
}

func waitForJob(jobID, apiKey string) error {
	url := fmt.Sprintf("%s/job/%s", baseURL, jobID)

	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Accept", "application/json")

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			bodyBytes, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("status check failed with status %d: %s", resp.StatusCode, string(bodyBytes))
		}

		var status JobStatus
		if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
			resp.Body.Close()
			return err
		}
		resp.Body.Close()

		if status.Status == "SUCCESS" || status.Status == "COMPLETED" {
			return nil
		}

		if status.Status == "ERROR" || status.Status == "FAILED" {
			return fmt.Errorf("job failed with status: %s", status.Status)
		}

		time.Sleep(pollDelay)
	}

	return fmt.Errorf("job timeout after %d retries", maxRetries)
}

func downloadMarkdown(jobID, apiKey string) (string, error) {
	url := fmt.Sprintf("%s/job/%s/result/markdown", baseURL, jobID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("download failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}
