use gemini_rust::Gemini;
use std::{env, error::Error, fs, path::Path};
use reqwest::Client;
use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize)]
struct Course {
    id: i32,
    title: String,
    color: String,
    icon: String,
    level: String,
    examination_body: String,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    
    // Gemini implimentation will abstract this somewhere to use a system prompt or instruction
    // with_system_prompt or with_system_instructions

    async fn ask_gemini() -> Result<(), Box<dyn Error>> {
        let api_key = env::var("GEMINI_API_KEY")?;
        let client = Gemini::new(api_key)?;

        let response = client
            .generate_content()
            .with_user_message("Hello, how are you?")
            .with_system_instruction("You are a top pyschologist and are empathetic to the last detail")
            .execute()
            .await?;

        println!("Response: {}", response.text());

        Ok(())
    }

   // ask_gemini().await?;

    let dir_path = Path::new("../courses");

    for entry in fs::read_dir(dir_path)? {
        let entry = entry?;
        let path = entry.path();

        println!("Reading Course {}", path.display());

        let syllabus = format!("{}/syllabus.json", path.display());
        let contents = fs::read_to_string(&syllabus)?;
        println!("The syllabus is {}", contents);
        let course: Course = serde_json::from_str(&contents)?;
        send_course(course).await?;
    }
    
   async fn send_course(course: Course) -> Result<(), Box<dyn Error>> {
        let client = Client::new();
        let response = client.post("http://localhost:8000/api/courses/")
            .json(&course)
            .send()
            .await?;

        println!("Response: {}", response.text().await?);
        Ok(())
    }

   send_course(Course {
        id: 4,
        title: "Introduction to Rust".to_string(),
        color: "blue".to_string(),
        icon: "rust-icon.png".to_string(),
        level: "beginner".to_string(),
        examination_body: "Rust University".to_string(),
    }).await?;



    Ok(())
}
