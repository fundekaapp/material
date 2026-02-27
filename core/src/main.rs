mod material_assistant;
mod material_processing;
mod upload_database;
use serde::Deserialize;

use crate::material_assistant::ask_gemini;
use crate::material_processing::{create_flashcards, create_lessons, refine_concepts};
use crate::upload_database::{Course, TopicMeta, send_concept, send_course, send_flashcards, send_lesson};
use std::{error::Error, fs, path::Path};

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    // ask_gemini().await?;

    let dir_path = Path::new("../courses");

    for entry in fs::read_dir(dir_path)? {
        let entry = entry?;
        let path = entry.path();

        println!("Reading Course {}", path.display());
        let syllabus = format!("{}/syllabus.json", path.display());
        let contents = fs::read_to_string(&syllabus)?;
        let course: Course = serde_json::from_str(&contents)?;
        // send_course(course).await?;
        // Enter course lessons
        let lesson_directory = format!("{}/topics", path.display());
        println!("lesson directory {}", lesson_directory);
        for lesson in fs::read_dir(lesson_directory)? {
            let lesson = lesson?;
            let lesson_path = lesson.path();

            if lesson_path.is_dir() {
                println!("Lesson path {}", lesson_path.display());
                let raw_concept_file = lesson_path.join("concepts_raw.md");
                println!("the raw concept path{}", raw_concept_file.display());
                
                
                let topic_file = lesson_path.join("topic.json");
                
                if topic_file.exists() {
                    let topic_string = fs::read_to_string(&topic_file)?;
                    let topic: TopicMeta = serde_json::from_str(&topic_string)?;
                    println!("Sending Concept : {}", topic_string);
                    send_flashcards(&lesson_path, topic.id(), course.id()).await?;

                    // changed the name to jsons just too break it for now
                     if lesson_path.join("concepts.jsons").exists() {
                     send_concept(&lesson_path, topic.id(), course.id()).await?;
                     }
                     else {
                         println!("concepts file does not exist {}", lesson_path.display())
                     }
                    // send_lesson(&lesson_path, course.id()).await?;

                }

                if raw_concept_file.exists() {
                    let raw_concepts = fs::read_to_string(&raw_concept_file)?;
                    // let refined_concepts = refine_concepts(raw_concepts).await?;
                    // let _ = fs::write(lesson_path.join("concepts.json"), &refined_concepts);

                    // let created_flashcards = create_flashcards(&refined_concepts).await?;
                    // fs::write(lesson_path.join("flashcards.json"), created_flashcards);

                    // let created_lesson = create_lessons(&refined_concepts).await?;
                    // fs::write(lesson_path.join("lesson.md"), created_lesson);
                } else {
                    println!("skipping {} no raw concepts here", lesson_path.display());
                }
            }
        }
    }

    Ok(())
}
