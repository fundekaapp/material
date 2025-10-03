package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"crypto/rand"
	"math/big"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Data structures for the improved schema
type Course struct {
	Title           string   `json:"title"`
	Color           string   `json:"color"`
	Icon            string   `json:"icon"`
	Lessons         []string `json:"lessons"`
	Level           string   `json:"level"`
	ExaminationBody string   `json:"examination_body"`
}

type Lesson struct {
	Icon          string   `json:"icon"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Content       string   `json:"content"`
	DeckIds       []string `json:"deck_ids"`
	QuizIds       []string `json:"quiz_ids"`
	Course        string   `json:"course"`
	Audio         string   `json:"audio,omitempty"`
	Concepts      []string `json:"concepts"`
	Order         int      `json:"order"`
	Prerequisites []string `json:"prerequisites"`
}

type TopicsData struct {
	Course  Course   `json:"course"`
	Lessons []Lesson `json:"lessons"`
}

// Question types for the centralized question bank
type Question struct {
	ID         string   `json:"id"`
	Type       string   `json:"type"`        // "mcq", "match", "oneword", "calculation", "flashcard"
	Question   string   `json:"question"`
	Answer     string   `json:"answer"`
	Options    []string `json:"options,omitempty"`    // For MCQ
	Concept    string   `json:"concept"`
	LessonID   string   `json:"lesson_id"`
	Source     string   `json:"source"`
	Difficulty string   `json:"difficulty"`
	Tags       []string `json:"tags"`
}

type Deck struct {
	ID                string   `json:"id"`
	Title             string   `json:"title"`
	LessonIDs         []string `json:"lesson_ids"`
	FlashcardQuestionIDs []string `json:"flashcard_question_ids"`
	StudyModeSettings map[string]interface{} `json:"study_mode_settings"`
}

type Quiz struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	LessonID     string   `json:"lesson_id"`
	QuestionIDs  []string `json:"question_ids"`
	TimeLimit    int      `json:"time_limit"`    // in minutes
	PassingScore int      `json:"passing_score"` // percentage
}

type QuestionsOutput struct {
	Questions []Question `json:"questions"`
	Decks     []Deck     `json:"decks"`
	Quizzes   []Quiz     `json:"quizzes"`
}

func generateQuestions() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println("❌ Please set the GEMINI_API_KEY environment variable.")
		return
	}

	// Paths
	topicsDir := "topics"
	questionsDir := "questions"

	// Ensure directories exist
	if err := os.MkdirAll(questionsDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create questions directory: %v\n", err)
		return
	}

	// Check if topics directory exists
	if _, err := os.Stat(topicsDir); os.IsNotExist(err) {
		fmt.Printf("❌ Topics directory '%s' does not exist. Run the topics processor first.\n", topicsDir)
		return
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		fmt.Printf("❌ Failed to create Gemini client: %v\n", err)
		return
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	// Process all topic JSON files
	files, err := ioutil.ReadDir(topicsDir)
	if err != nil {
		fmt.Printf("❌ Failed to read topics directory: %v\n", err)
		return
	}

	processedCount := 0
	skippedCount := 0

	for _, f := range files {
		if filepath.Ext(f.Name()) != ".json" {
			continue
		}

		topicsPath := filepath.Join(topicsDir, f.Name())
		questionsName := strings.TrimSuffix(f.Name(), ".json") + "_questions.json"
		questionsPath := filepath.Join(questionsDir, questionsName)

		// Skip if already exists
		if _, err := os.Stat(questionsPath); err == nil {
			fmt.Printf("✅ Skipping %s (already has questions)\n", f.Name())
			skippedCount++
			continue
		}

		fmt.Printf("🔄 Generating questions for %s...\n", f.Name())

		// Read topics JSON
		content, err := ioutil.ReadFile(topicsPath)
		if err != nil {
			fmt.Printf("❌ Failed to read topics file '%s': %v\n", f.Name(), err)
			continue
		}

		// Parse topics data
		var topicsData TopicsData
		if err := json.Unmarshal(content, &topicsData); err != nil {
			fmt.Printf("❌ Failed to parse topics JSON '%s': %v\n", f.Name(), err)
			continue
		}

		// Generate questions for all lessons
		questionsOutput := QuestionsOutput{
			Questions: []Question{},
			Decks:     []Deck{},
			Quizzes:   []Quiz{},
		}

		for _, lesson := range topicsData.Lessons {
			// Generate questions for this lesson
			lessonQuestions, err := generateQuestionsForLesson(ctx, model, lesson, topicsData.Course)
			if err != nil {
				fmt.Printf("⚠️  Failed to generate questions for lesson '%s': %v\n", lesson.Title, err)
				continue
			}

			questionsOutput.Questions = append(questionsOutput.Questions, lessonQuestions...)

			// Create flashcard deck for this lesson
			flashcardIDs := []string{}
			for _, q := range lessonQuestions {
				if q.Type == "flashcard" {
					flashcardIDs = append(flashcardIDs, q.ID)
				}
			}

			if len(flashcardIDs) > 0 {
				deck := Deck{
					ID:                   generateID(),
					Title:                fmt.Sprintf("%s - Flashcards", lesson.Title),
					LessonIDs:           []string{fmt.Sprintf("lesson_%d", lesson.Order)},
					FlashcardQuestionIDs: flashcardIDs,
					StudyModeSettings: map[string]interface{}{
						"show_answer_delay": 3,
						"shuffle_cards":     true,
						"repeat_incorrect":  true,
					},
				}
				questionsOutput.Decks = append(questionsOutput.Decks, deck)
			}

			// Create quiz for this lesson
			quizQuestionIDs := []string{}
			for _, q := range lessonQuestions {
				if q.Type != "flashcard" {
					quizQuestionIDs = append(quizQuestionIDs, q.ID)
				}
			}

			if len(quizQuestionIDs) > 0 {
				quiz := Quiz{
					ID:           generateID(),
					Title:        fmt.Sprintf("%s - Practice Quiz", lesson.Title),
					LessonID:     fmt.Sprintf("lesson_%d", lesson.Order),
					QuestionIDs:  quizQuestionIDs,
					TimeLimit:    30, // 30 minutes default
					PassingScore: 70, // 70% default
				}
				questionsOutput.Quizzes = append(questionsOutput.Quizzes, quiz)
			}

			fmt.Printf("  ✅ Generated %d questions for '%s'\n", len(lessonQuestions), lesson.Title)
		}

		// Write questions JSON file
		formattedJSON, err := json.MarshalIndent(questionsOutput, "", "  ")
		if err != nil {
			fmt.Printf("❌ Failed to format questions JSON for '%s': %v\n", f.Name(), err)
			continue
		}

		if err := os.WriteFile(questionsPath, formattedJSON, 0644); err != nil {
			fmt.Printf("❌ Failed to write questions file '%s': %v\n", questionsName, err)
			continue
		}

		fmt.Printf("✅ Saved %s (%d questions, %d decks, %d quizzes)\n", 
			questionsName, len(questionsOutput.Questions), len(questionsOutput.Decks), len(questionsOutput.Quizzes))
		processedCount++

		// Rate limiting
		time.Sleep(10 * time.Second)
	}

	// Summary
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("🎉 Questions generation complete!\n")
	fmt.Printf("📊 Files processed: %d\n", processedCount)
	fmt.Printf("⏭️  Files skipped: %d\n", skippedCount)
	fmt.Printf("📁 Output directory: %s\n", questionsDir)
}

func generateQuestionsForLesson(ctx context.Context, model *genai.GenerativeModel, lesson Lesson, course Course) ([]Question, error) {
	prompt := fmt.Sprintf(`You are an educational content expert. Generate questions and flashcards for the following lesson:

LESSON: %s
COURSE: %s (%s)
DESCRIPTION: %s
CONCEPTS: %s

Generate a variety of questions for each concept:

1. FLASHCARDS (3-4 per concept):
   - Simple definition cards: "What is [concept]?" -> "Definition"
   - Key fact cards: "What does [concept] do?" -> "Function/Purpose"
   - Example cards: "Give an example of [concept]" -> "Specific example"

2. QUIZ QUESTIONS (2-3 per concept):
   - Multiple choice questions (4 options)
   - One word/short answer questions
   - Calculation questions (if applicable)

For each question, determine appropriate difficulty: "easy", "medium", "hard"

Return ONLY valid JSON in this exact format:
[
  {
    "id": "unique_id",
    "type": "flashcard",
    "question": "What is photosynthesis?",
    "answer": "The process by which plants convert light energy into chemical energy",
    "concept": "photosynthesis",
    "difficulty": "easy",
    "tags": ["definition", "biology", "plants"]
  },
  {
    "id": "unique_id",
    "type": "mcq",
    "question": "Which organelle is responsible for photosynthesis?",
    "answer": "Chloroplasts",
    "options": ["Mitochondria", "Chloroplasts", "Nucleus", "Ribosomes"],
    "concept": "photosynthesis",
    "difficulty": "medium",
    "tags": ["organelles", "cell structure"]
  },
  {
    "id": "unique_id",
    "type": "oneword",
    "question": "What gas is released during photosynthesis?",
    "answer": "Oxygen",
    "concept": "photosynthesis",
    "difficulty": "easy",
    "tags": ["gas exchange", "products"]
  }
]

Generate questions for concepts: %s`,
		lesson.Title, course.Title, course.Level, lesson.Description, 
		strings.Join(lesson.Concepts, ", "), strings.Join(lesson.Concepts, ", "))

	// Send to Gemini
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	// Collect response
	var out strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content == nil {
			continue
		}
		for _, part := range cand.Content.Parts {
			out.WriteString(fmt.Sprint(part))
		}
	}

	responseText := out.String()
	if responseText == "" {
		return nil, fmt.Errorf("empty response from Gemini")
	}

	// Clean JSON response
	cleanedJSON := cleanJSONResponse(responseText)

	// Parse questions
	var questions []Question
	if err := json.Unmarshal([]byte(cleanedJSON), &questions); err != nil {
		return nil, fmt.Errorf("failed to parse questions JSON: %v\nResponse: %s", err, cleanedJSON)
	}

	// Post-process questions
	for i := range questions {
		if questions[i].ID == "" {
			questions[i].ID = generateID()
		}
		questions[i].LessonID = fmt.Sprintf("lesson_%d", lesson.Order)
		questions[i].Source = "generated"
		
		if questions[i].Tags == nil {
			questions[i].Tags = []string{}
		}
	}

	return questions, nil
}

// cleanJSONResponse removes common formatting issues from AI responses
func cleanJSONResponse(response string) string {
	// Remove markdown code blocks
	re := regexp.MustCompile("```(?:json)?\n?(.*?)\n?```")
	matches := re.FindStringSubmatch(response)
	if len(matches) > 1 {
		response = matches[1]
	}

	// Trim whitespace
	response = strings.TrimSpace(response)

	// Remove any leading/trailing non-JSON content
	startIdx := strings.Index(response, "[")
	endIdx := strings.LastIndex(response, "]")
	
	if startIdx >= 0 && endIdx >= 0 && endIdx > startIdx {
		response = response[startIdx : endIdx+1]
	}

	return response
}

// generateID creates a simple random ID
func generateID() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 8)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		result[i] = chars[num.Int64()]
	}
	return string(result)
}
