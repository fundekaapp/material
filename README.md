# Fundeka Material Processor

> Content pipeline and structured learning materials for the [Fundeka](https://fundeka.app) learning platform.

## Overview

This repository serves as the single source of truth for all educational content in Fundeka. It combines structured course materials with an automated processing pipeline (`core/`) that uses the Gemini AI API to generate and enrich content — flashcards, lessons, topic metadata, and more.

---

## Repository Structure

```
fundeka-materials/
├── README.md
├── core/                  # Rust-based content processor (Gemini-powered)
├── courses/               # Structured course content
├── exam_papers/           # Past exam PDFs by examination body
└── sylubi/                # Official syllabus PDFs by subject
```

### `courses/`

Each course follows the naming convention: **`examinationbody_subject_level`**

```
courses/
└── zimsec_mathematics_olevel/
    ├── syllabus.json          # Course metadata (description, color, icon)
    └── topics/
        └── algebra/
            ├── topic.json         # Topic metadata
            ├── lesson.md          # Lesson content (Markdown)
            ├── lesson.pdf         # Lesson content (PDF)
            ├── flashcards.json    # Flashcard Q&A pairs
            ├── assets/            # Images and other media
            └── videos/            # Video resources
```

#### `syllabus.json` schema

```json
{
  "title": "Mathematics",
  "examination_body": "ZIMSEC",
  "level": "O Level",
  "description": "A brief description of the course.",
  "color": "#4F46E5",
  "icon": "calculator"
}
```

#### `topic.json` schema

```json
{
  "title": "Algebra",
  "description": "Introduction to algebraic expressions and equations.",
  "order": 1
}
```

#### `flashcards.json` schema

```json
[
  {
    "question": "What is the quadratic formula?",
    "answer": "x = (-b ± √(b²-4ac)) / 2a"
  }
]
```

### `exam_papers/`

PDF past papers organised by examination body and subject.

### `sylubi/`

Official syllabus PDFs for reference during content generation.

### `core/`

A Rust program that automates content processing across the repository using the [Gemini API](https://ai.google.dev/). It walks the `courses/` directory, reads each `syllabus.json`, and will eventually orchestrate AI-driven generation of lessons, flashcards, and topic content.

**Status:** Active development — currently boilerplate with Gemini client integration scaffolded.

**Prerequisites:**
- Rust (stable)
- A valid `GEMINI_API_KEY` environment variable

**Running the processor:**

```bash
cd core
export GEMINI_API_KEY=your_key_here
cargo run
```

---

## Adding a New Course

1. Create a directory under `courses/` using the `examinationbody_subject_level` convention (lowercase, underscores).
2. Add a `syllabus.json` with course metadata.
3. Create a `topics/` directory and add one subdirectory per topic.
4. Populate each topic directory with `topic.json`, `lesson.md`, `lesson.pdf`, `flashcards.json`, and any `assets/` or `videos/`.

---

## Contributing

Content contributions (corrections, additional topics, new exam papers) are welcome. Please keep filenames lowercase with underscores and validate JSON files before opening a pull request.

---

## License

Content and code are proprietary to Fundeka unless otherwise stated.
