# Material Processor for Fundeka

This repository contains tools and resources for collecting, parsing, and organizing educational material for **Fundeka**. The primary focus is on converting official syllabus documents and exam papers into structured markdown that can later be added to the database and used to generate learning content.

---

## Overview

Although the contents of this repository are proprietary, they may be useful for both **students** and **developers**:

* **Students** benefit from access to processed syllabus content, official documents, and relevant exam papers.
* **Developers** can use the parsed syllabus, material, and questions to:

  * Build their own content pipelines.
  * Create content collections for schools aligned with the official syllabus.

---

## Current Structure

* **`parser/`** – Contains the raw syllabus PDFs as well as the parsed markdown versions.
* **`markdown/`** – Contains cleaned and processed markdown files.

---

## Usage

1. Add the latest syllabus documents to the **`parser/pdf/`** folder.
2. Run the parser:

   ```bash
   go run parse.go
   ```
3. Clean the generated markdown:

   ```bash
   lua clean_markdown.lua
   ```

*(Note: This process will be streamlined in the future, potentially automated via GitHub Actions.)*

---
# Course Structure
 - ExaminationBody_Subject_Grade/
    - Sylubus.md
    - examination_dates.json
    - topics/
        - lesson.md
        - flashcards.json
        - questions.json
        - videos/
        - audios/
---

## Roadmap

The following features and improvements are planned:

* **Content Structuring**

  * Extract individual grades, topics, concepts, and outcomes from syllabus markdown.
  * Push structured concepts into the database.

* **Lesson Material Generation**

  * Automatically generate lesson material, flashcards, and quizzes from concepts.

* **Past Paper Integration**

  * Parse past papers and extract questions/answers.
  * Link quiz questions to their original sources for reference.

* **Media Generation**

  * Generate video and podcast-style audio scripts for concepts and lessons. check out [Vilyrean](https://github.com/mchiwundura/vilyrian)
  * Create explainer shorts and stitch them into long-form content (e.g., YouTube).
  * Produce question explanation videos for past exam papers.

---

## Contributing

Contributions are welcome! You can:

* Improve the parser.
* Update syllabus documents with the latest versions.
* Suggest or implement automation (e.g., GitHub Actions).

Please ensure contributions align with the repository’s purpose of supporting structured educational material for Fundeka.

---

## License

This project is proprietary. Unauthorized redistribution of syllabus documents or exam papers may be restricted by copyright.

---

