use gemini_rust::Gemini;
use std::{env, error::Error};

pub async fn refine_concepts(raw_concepts: String) -> Result<String, Box<dyn Error>> {
    
    let prompt = format!("The raw concepts: {}", raw_concepts);

        let api_key = env::var("GEMINI_API_KEY")?;
        let client = Gemini::new(api_key)?;

        let response = client
            .generate_content()
            .with_user_message(prompt)
            .with_system_instruction("You will receive concepts for a topic written raw, you should return these concepets in json formated with an id,title and body for the concept the first part of the id is provided at the top therefore for the concepts add the last part starting with C and increment the numbers DBE-PS12-MI-C01, just give the raw json without markdown formating and just start off with the object containing the concepts array ")
            .execute()
            .await?;
    
    Ok(response.text())
}

pub async fn create_flashcards(concepts: &String) -> Result<String, Box<dyn Error>> {
    
    let prompt = format!("The raw concepts: {}", concepts);

        let api_key = env::var("GEMINI_API_KEY")?;
        let client = Gemini::new(api_key)?;

        let response = client
            .generate_content()
            .with_user_message(prompt)
            .with_system_instruction("You will receive concepts in json for these concepts you should create flashcards from these concepts and output them in json without the outer markdonw formatting starting off with the object containing a flashcards array the flashcard items in the array should have the id formated as  the concept_id but difference is the end part to have the FC prefix for example DBE-PS12-MI-FC01, concept_id, front and back these flashcards are to assist memorize the concept and learn it so the front is a question on the concept and the back is the answer for example if the concept is defining momentum the question can be define momentum and the back will be the definition of momemntum")
            .execute()
            .await?;
    
    Ok(response.text())
}

pub async fn create_lessons(concepts: &String) -> Result<String, Box<dyn Error>> {
    
    let prompt = format!("The raw concepts: {}", concepts);

        let api_key = env::var("GEMINI_API_KEY")?;
        let client = Gemini::new(api_key)?;

        let response = client
            .generate_content()
            .with_user_message(prompt)
            .with_system_instruction("you will receive concepts in json from these concepts you will create a lesson in markdown to teach these concepts to a student for things requiring memoriztoin you can create mnemonics for the student and examples based on their interests, the student likes soccer and is a fan of real madrid he also likes rap music like likes from kanye west and enjoys christopher nolan movies do not over do these but good in specific examples and mnemonics, do not mention the actual concept ids or number just title the sections of the lesson accordingly to the title of the concept")
            .execute()
            .await?;
    
    Ok(response.text())
}

