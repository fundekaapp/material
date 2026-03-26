use gemini_rust::{ContentBuilder, Gemini};
use std::{env, error::Error, fs, path::Path};

pub async fn refine_concepts(raw_concepts: String) -> Result<String, Box<dyn Error>> {
    let prompt = format!("The raw concepts: {}", raw_concepts);

    let api_key = env::var("GEMINI_API_KEY")?;
    let client = Gemini::new(api_key)?;

    let response = client
            .generate_content()
            .with_user_message(prompt)
            .with_system_instruction("You will receive concepts for a topic written raw, you should return these concepets in json formated with an id,title and body for the concept the first part of the id is provided at the top therefore for the concepts add the last part starting with C and increment the numbers DBE-PS12-MI-C01, and a a children array where you will add child concepts to the current concept if they are related in that order, just give the raw json without markdown formating and just start off with the object containing the concepts array ")
            .execute()
            .await?;

    Ok(response.text())
}

pub async fn create_topic(raw_concepts: String) -> Result<String, Box<dyn Error>> {
    let instructions_path = Path::new("../instructions/create_topic.md");
    let instructions = fs::read_to_string(instructions_path)?;

    let prompt = format!("The raw concepts: {}", raw_concepts);

    let api_key = env::var("GEMINI_API_KEY")?;
    let client = Gemini::new(api_key)?;

    let response = client
        .generate_content()
        .with_user_message(prompt)
        .with_system_instruction(instructions)
        .with_response_mime_type("application/json")
        .with_temperature(0.1)
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

pub async fn list_topics_from_pdf(
    pdf_path: &Path,
    course_path: &String,
) -> Result<String, Box<dyn Error>> {
    let pdf_bytes = fs::read(pdf_path)?;

    let api_key = env::var("GEMINI_API_KEY")?;
    let client = Gemini::new(api_key)?;

    let instructions_path = Path::new("../instructions/create_topics_from_pdf.md");
    let instructions = fs::read_to_string(instructions_path)?;

    let file_handle = client
        .create_file(pdf_bytes)
        .display_name(pdf_path.file_name().unwrap().to_string_lossy())
        .with_mime_type("application/pdf".parse().unwrap())
        .upload()
        .await?;

    let response = client
        .generate_content()
        .with_system_instruction(instructions)
        .with_response_mime_type("application/json")
        .with_user_message_and_file(course_path, &file_handle)?
        .execute()
        .await?;

    let _ = file_handle.delete().await;

    Ok(response.text())
}

pub async fn create_topic_concepts_from_pdf(
    pdf_path: &Path,
    topic_info: &String,
) -> Result<String, Box<dyn Error>> {
    let pdf_bytes = fs::read(pdf_path)?;

    let api_key = env::var("GEMINI_API_KEY")?;
    let client = Gemini::new(api_key)?;

    let instructions_path = Path::new("../instructions/create_topic_concepts_from_pdf.md");
    let instructions = fs::read_to_string(instructions_path)?;

    let file_handle = client
        .create_file(pdf_bytes)
        .display_name(pdf_path.file_name().unwrap().to_string_lossy())
        .with_mime_type("application/pdf".parse().unwrap())
        .upload()
        .await?;

    let response = client
        .generate_content()
        .with_system_instruction(instructions)
        .with_response_mime_type("application/json")
        .with_user_message_and_file(topic_info, &file_handle)?
        .execute()
        .await?;

    let _ = file_handle.delete().await;

    Ok(response.text())
}
