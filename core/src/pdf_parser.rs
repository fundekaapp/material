use pdf_oxide::PdfDocument;
use pdf_oxide::pipeline::OutputConverter;
use pdf_oxide::pipeline::{TextPipeline, TextPipelineConfig};
use pdf_oxide::pipeline::converters::MarkdownOutputConverter;

// Open a PDF
pub fn parse_pdf() -> Result<(), Box<dyn std::error::Error>> {

let mut doc = PdfDocument::open("test.pdf")?;

// Extract text with reading order (multi-column support)
 let spans = doc.extract_spans(10)?;
 let config = TextPipelineConfig::default();
 let pipeline = TextPipeline::with_config(config.clone());
 let ordered_spans = pipeline.process(spans, Default::default())?;

//  Convert to Markdown
 let converter = MarkdownOutputConverter::new();
 let markdown = converter.convert(&ordered_spans, &config)?;
 println!("{}", markdown);


Ok(())
}
