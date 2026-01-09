Document Processing System Implementation
This plan outlines the implementation of document processing capabilities in the ERP6 system, allowing the AI component to read and learn from PDF and DOC/DOCX files.

User Review Required
IMPORTANT

New Dependencies Required

We need to add github.com/nguyenthenguyen/docx for DOC/DOCX text extraction
The PDF library github.com/ledongthuc/pdf is already installed
NOTE

Document Storage Strategy

Documents will be stored in the database with their extracted text content
Large documents may impact database size - consider adding a file size limit (e.g., 10MB)
Document content will be injected into AI prompts, which may increase token usage
Proposed Changes
Backend - Document Processing Component
[NEW] 
document.go
Create a new model for storing document metadata and extracted content:

DocumentID (primary key)
UserID (who uploaded it)
FileName (original filename)
FilePath (storage path)
FileType (pdf, doc, docx)
FileSize (bytes)
ExtractedText (full text content)
Summary
 (optional AI-generated summary)
CreatedAt, UpdatedAt (timestamps)
Backend - Component Implementation
[NEW] 
comp_document.go
Create a new workflow component for document processing with the following capabilities:

Action: upload_and_extract - Upload a document file and extract its text content

Parameters: file_field, action (upload_and_extract, extract_only, list, get)
Extract text from PDF using github.com/ledongthuc/pdf
Extract text from DOC/DOCX using github.com/nguyenthenguyen/docx
Store document metadata and extracted text in the documents table
Return document ID and extracted text in workflow result
Inject document content into ctx.Extras for use by other components
Action: list - List all documents for the current user

Returns array of document metadata
Action: 
get
 - Retrieve a specific document's content

Parameters: document_id
Returns full document metadata and extracted text
Action: delete - Delete a document

Parameters: document_id
Removes file from disk and database record
Backend - AI Component Enhancement
[MODIFY] 
comp_ai.go
Enhance the AI component to accept and use document context:

Add new parameter: document_ids (comma-separated list of document IDs to include as context)
Add new parameter: use_all_documents (boolean, include all user's documents)
In 
delegateToLLM
 function, fetch document content from database
Inject document content into the system prompt before sending to LLM
Format document context as: "Available Documents:\n\n[Document: filename.pdf]\n{extracted_text}\n\n"
Update conversation state to track which documents are being referenced
Database Migration
[NEW] 
create_documents_table.sql
Create SQL migration for the documents table:

CREATE TABLE IF NOT EXISTS documents (
    documentid SERIAL PRIMARY KEY,
    userid INT NOT NULL,
    filename VARCHAR(255) NOT NULL,
    filepath VARCHAR(500) NOT NULL,
    filetype VARCHAR(10) NOT NULL,
    filesize BIGINT NOT NULL,
    extractedtext TEXT,
    summary TEXT,
    createdat TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedat TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_documents_userid ON documents(userid);
CREATE INDEX idx_documents_createdat ON documents(createdat);
Backend - Component Registry
[MODIFY] 
registry.go
No changes needed - the comp_document component will auto-register via its 
init()
 function.

Frontend - Component Configuration
[MODIFY] Component Configuration (Frontend)
Add comp_document to the workflow designer component library:

Component name: "Document Processor"
Category: "Data Processing"
Icon: Document/File icon
Parameters:
action (select: upload_and_extract, list, get, delete)
file_field (text input, required for upload_and_extract)
document_id (text input, required for get/delete)
[MODIFY] AI Component Configuration (Frontend)
Add new parameters to the AI component:

document_ids (text input, placeholder: "1,2,3 or leave empty")
use_all_documents (checkbox, default: false)
Verification Plan
Manual Testing
Since this is a new feature with workflow integration, manual testing is the most appropriate approach:

Test 1: Document Upload and Extraction
Start the backend server: go run main.go from 
/c:/lara/www/erp6-be-golang
Create a new workflow in the frontend with the following nodes:

Node 1: comp_document with action=upload_and_extract, file_field=document

Node 2: comp_table to save the document ID to a test table
Upload a test PDF file (create a simple PDF with text like "This is a test document for AI learning")
Verify the document is saved in the documents table
Verify the extractedtext field contains the correct text from the PDF
Check the console logs for any errors

Test 2: AI with Document Context
Create a workflow with:
Node 1: comp_ai with document_ids set to the ID from Test 1
Configure the AI to answer: "What does the document say?"
Execute the workflow
Verify the AI response references the content from the uploaded document
Check that the AI's response accurately reflects the document content

Test 3: List and Get Documents
Create a workflow with comp_document action=list
Verify it returns all documents for the current user
Create another workflow with comp_document action=
get
 and a specific document_id
Verify it returns the full document content
Test 4: DOCX File Support
Create a test DOCX file with text content
Upload it using the same workflow from Test 1
Verify the text is correctly extracted from the DOCX file
Database Verification
After running tests, manually check the documents table:
SELECT * FROM documents ORDER BY createdat DESC LIMIT 5;
Verify all fields are populated correctly
Check that extractedtext contains readable content
Error Handling Tests
Test uploading an unsupported file type (e.g., .txt, .jpg)
Test uploading a very large file (>10MB if limit is implemented)
Test accessing a document that doesn't exist
Test accessing another user's document (should fail)