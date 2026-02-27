use reqwest::Client;
use serde::{Deserialize, Serialize};
use std::{error::Error, fs, path::PathBuf};

#[derive(Serialize, Deserialize)]
pub struct Course {
    id: i32,
    title: String,
    color: String,
    icon: String,
    level: String,
    examination_body: String,
}

impl Course {
    pub fn id(&self) -> i32 {
        self.id
    }
}

#[derive(Serialize, Deserialize)]
pub struct Lesson {
    id: String,
    icon: String,
    title: String,
    description: String,
    content: String,
    audio: String,
    order: i32,
    course: i32,
}

#[derive(Deserialize)]
pub struct TopicMeta {
    id: String,
    title: String,
    description: String,
    order: i32,
    icon: String,
}

impl TopicMeta {
    pub fn id(&self) -> String {
        self.id.clone()
    }
}

#[derive(Serialize, Deserialize)]
pub struct Concept {
    id: String,
    title: String,
    outcome: String,
    order: i32,
    lesson: String,
    // parent: String,
    difficulty: i32,
}

#[derive(Serialize, Deserialize)]
pub struct ConceptMeta {
    id: String,
    title: String,
    body: String,
}

#[derive(Serialize, Deserialize)]
pub struct ConceptArray {
    concepts: Vec<ConceptMeta>,
}

#[derive(Serialize, Deserialize)]
pub struct Flashcard {
    id: String,
    front: String,
    back: String,
    lesson: String,
    concept: String,
}

#[derive(Serialize, Deserialize)]
pub struct FlashcardMeta {
    id: String,
    front: String,
    back: String,
    concept_id: String,
}


#[derive(Serialize, Deserialize)]
pub struct FlashcardArray {
 flashcards : Vec<FlashcardMeta>
}

pub async fn send_course(course: Course) -> Result<(), Box<dyn Error>> {
    let client = Client::new();
    let response = client
        .post("http://localhost:8000/api/courses/")
        .json(&course)
        .send()
        .await?;

    println!("Response: {}", response.text().await?);
    Ok(())
}

pub async fn send_lesson(file_path: &PathBuf, course_id: i32) -> Result<(), Box<dyn Error>> {
    let content = fs::read_to_string(file_path.join("lesson.md"))?;
    let audio_path = file_path.join("flashcards.json");
    let audio = audio_path.to_string_lossy().to_string();

    let topic_string = fs::read_to_string(file_path.join("topic.json"))?;
    let topic: TopicMeta = serde_json::from_str(&topic_string)?;

    let lesson = Lesson {
        id: format!("{}{}", course_id, topic.id),
        icon: topic.icon,
        title: topic.title,
        description: topic.description,
        content,
        audio,
        order: topic.order,
        course: course_id,
    };

    let client = Client::new();
    let response = client
        .post("http://localhost:8000/api/lessons/")
        .json(&lesson)
        .send()
        .await?;

    println!("Response: {}", response.text().await?);
    Ok(())
}

pub async fn send_concept(file_path: &PathBuf, lesson_id: String, course_id: i32) -> Result<(), Box<dyn Error>> {
    let concepts_string = fs::read_to_string(file_path.join("concepts.json"))?;
    let concepts: ConceptArray = serde_json::from_str(&concepts_string)?;

    for (index, concept) in concepts.concepts.iter().enumerate() {
        let concept = Concept {
            id: concept.id.clone(),
            title: concept.title.clone(),
            outcome: concept.body.clone(),
            order: index as i32,
            difficulty: 0,
            lesson: format!("{}{}", course_id, lesson_id.clone()),
       //     parent: concept.id.clone(),
        };

        let client = Client::new();
        let response = client
            .post("http://localhost:8000/api/concepts/")
            .json(&concept)
            .send()
            .await?;

        println!("Response: {}", response.text().await?);
    }
    Ok(())
}

pub async fn send_flashcards(file_path: &PathBuf, lesson_id: String, course_id: i32) -> Result<(), Box<dyn Error>> {
    let flashcards_string = fs::read_to_string(file_path.join("flashcards.json"))?;
    let flashcards: FlashcardArray = serde_json::from_str(&flashcards_string)?;
 
    for flashcard in flashcards.flashcards.iter() {
        
        let flashcard = Flashcard {
            id: flashcard.id.clone(),
            front: flashcard.front.clone(),
            back: flashcard.back.clone(),
            lesson: format!("{}{}", course_id, lesson_id.clone()),
            concept: flashcard.concept_id.clone(),
        };

        let client = Client::new();
        let response = client
            .post("http://localhost:8000/api/flashcards/")
            .json(&flashcard)
            .send()
            .await?;

        println!("Response: {}", response.text().await?);
    }

    Ok(())
}




