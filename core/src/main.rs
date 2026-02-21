use gemini_rust::Gemini;
use std::{env, error::Error, fs, path::Path};

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
    ask_gemini().await?;

    let dir_path = Path::new("../courses");

    for entry in fs::read_dir(dir_path)? {
        let entry = entry?;
        let path = entry.path();

        println!("Reading Course {}", path.display());

        let syllabus = format!("{}/syllabus.json", path.display());
        let contents = fs::read_to_string(&syllabus)?;
        println!("The syllabus is {}", contents)
    }
    Ok(())
}
