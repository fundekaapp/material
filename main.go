package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"github.com/ledongthuc/pdf"
)

const chunkSize = 3500 // characters per chunk, safe size for Gemini

func main() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set the GEMINI_API_KEY environment variable.")
		return
	}

	// Paths
	pdfDir := "pdf"
	mdDir := "markdown"
	instructionsFile := "instructions.md"

	// Ensure markdown dir exists
	if err := os.MkdirAll(mdDir, 0755); err != nil {
		fmt.Println("Failed to create markdown dir:", err)
		return
	}

	// Load instructions
	instructions, err := ioutil.ReadFile(instructionsFile)
	if err != nil {
		fmt.Println("Failed to read instructions.md:", err)
		return
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		fmt.Println("Failed to create Gemini client:", err)
		return
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	files, _ := ioutil.ReadDir(pdfDir)
	for _, f := range files {
		if filepath.Ext(f.Name()) != ".pdf" {
			continue
		}

		pdfPath := filepath.Join(pdfDir, f.Name())
		mdName := strings.TrimSuffix(f.Name(), ".pdf") + ".md"
		mdPath := filepath.Join(mdDir, mdName)

		// Skip if already exists
		if _, err := os.Stat(mdPath); err == nil {
			fmt.Printf("✅ Skipping %s (already converted)\n", f.Name())
			continue
		}

		fmt.Printf("🔄 Processing %s...\n", f.Name())

		// Extract PDF text
		content, err := extractPDFText(pdfPath)
		if err != nil {
			fmt.Println("❌ Failed to extract text:", err)
			continue
		}

		// Chunk the content
		chunks := chunkText(content, chunkSize)
		fmt.Printf("📖 Split into %d chunks\n", len(chunks))

		var finalOutput strings.Builder
		for i, chunk := range chunks {
			fmt.Printf("➡️  Sending chunk %d/%d...\n", i+1, len(chunks))

			prompt := fmt.Sprintf(
				"Instructions:\n%s\n\nPDF chunk:\n%s\n\n---\nExtract and list the **individual concepts** that learners need to understand from this content. Group them under their **topics**, and structure the output as Markdown.",
				string(instructions), chunk,
			)

			resp, err := model.GenerateContent(ctx, genai.Text(prompt))
			if err != nil {
				fmt.Println("❌ Gemini error:", err)
				continue
			}

			for _, cand := range resp.Candidates {
				for _, part := range cand.Content.Parts {
					finalOutput.WriteString(fmt.Sprint(part))
				}
			}
			finalOutput.WriteString("\n\n")

			// Wait between requests
			time.Sleep(20 * time.Second)
		}

		// Write markdown file
		if err := os.WriteFile(mdPath, []byte(finalOutput.String()), 0644); err != nil {
			fmt.Println("❌ Failed to write file:", err)
			continue
		}

		fmt.Printf("✅ Saved %s\n", mdName)
	}

	fmt.Println("🎉 All PDFs processed.")

	// Generate topics
	fmt.Println("▶️ Starting topic generation...")
	generateTopics()

	// Generate questions
	fmt.Println("▶️ Starting question generation...")
	generateQuestions()
}

// extractPDFText extracts plain text from a PDF
func extractPDFText(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var buf strings.Builder
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	_, err = io.Copy(&buf, b)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// chunkText splits text into ~chunkSize character chunks
func chunkText(s string, size int) []string {
	var chunks []string
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanRunes)

	var buf strings.Builder
	for scanner.Scan() {
		buf.WriteString(scanner.Text())
		if buf.Len() >= size {
			chunks = append(chunks, buf.String())
			buf.Reset()
		}
	}
	if buf.Len() > 0 {
		chunks = append(chunks, buf.String())
	}
	return chunks
}

