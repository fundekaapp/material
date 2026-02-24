use gemini_rust::Gemini;
use std::{env, error::Error};

pub async fn ask_gemini() -> Result<(), Box<dyn Error>> {
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
