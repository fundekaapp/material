use reqwest::Client;
use serde::{Deserialize, Serialize};
use std::{error::Error};

#[derive(Serialize, Deserialize)]
pub struct Course {
    id: i32,
    title: String,
    color: String,
    icon: String,
    level: String,
    examination_body: String,
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
