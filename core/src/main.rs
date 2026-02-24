mod material_assistant;
mod material_processing;
mod upload_database;
use crate::material_assistant::ask_gemini;
use crate::material_processing::refine_concepts;
use crate::upload_database::{Course, send_course};
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

                if raw_concept_file.exists() {
                    let raw_concepts = fs::read_to_string(&raw_concept_file)?;
                    println!("the raw concepts {}", raw_concepts);
                    refine_concepts(raw_concepts);
                } else {
                    println!("skipping {} no raw concepts here", lesson_path.display());
                }
            }
        }
    }


    Ok(())
}
