package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func generateTopics() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set the GEMINI_API_KEY environment variable.")
		return
	}

	// Paths
	mdDir := "markdown"
	outDir := "topics"

	// Ensure topics dir exists
	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Println("Failed to create topics dir:", err)
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

	// Walk through Markdown files
	files, _ := ioutil.ReadDir(mdDir)
	for _, f := range files {
		if filepath.Ext(f.Name()) != ".md" {
			continue
		}

		mdPath := filepath.Join(mdDir, f.Name())
		jsonName := strings.TrimSuffix(f.Name(), ".md") + ".json"
		jsonPath := filepath.Join(outDir, jsonName)

		// Skip if already exists
		if _, err := os.Stat(jsonPath); err == nil {
			fmt.Printf("✅ Skipping %s (already has JSON)\n", f.Name())
			continue
		}

		fmt.Printf("🔄 Processing %s...\n", f.Name())

		// Read Markdown
		content, err := ioutil.ReadFile(mdPath)
		if err != nil {
			fmt.Println("❌ Failed to read markdown:", err)
			continue
		}

		// Prompt
		prompt := fmt.Sprintf(`You are given a syllabus in Markdown. 
Extract the list of **topics** and classify each with its **level** (e.g., "Form 5", "Form 6", "Grade 10"). 
Return only valid JSON in the format:

{
  "course": {title, icon, color, level, examination_body, lessons}
  "topics": [
    { "lesson": "<topic name>", "level": "<level>" },
    ...
  ]
}

Syllabus content:
%s`, string(content))

		// Send to Gemini
		resp, err := model.GenerateContent(ctx, genai.Text(prompt))
		if err != nil {
			fmt.Println("❌ Gemini error:", err)
			continue
		}

		// Collect response
		var out strings.Builder
		for _, cand := range resp.Candidates {
			for _, part := range cand.Content.Parts {
				out.WriteString(fmt.Sprint(part))
			}
		}

		// Write JSON file
		if err := os.WriteFile(jsonPath, []byte(out.String()), 0644); err != nil {
			fmt.Println("❌ Failed to write JSON file:", err)
			continue
		}

		fmt.Printf("✅ Saved %s\n", jsonName)

		// Wait to avoid hammering API
		time.Sleep(20 * time.Second)
	}

	fmt.Println("🎉 All Markdown syllabuses processed into JSON topics.")
}

