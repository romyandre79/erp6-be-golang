Document Processing System - Implementation Walkthrough
Overview
Successfully implemented a comprehensive document processing system for the ERP6 application, enabling the AI component to read and learn from PDF and DOC/DOCX files.

Changes Made
1. Dependencies
Installed Library:

github.com/nguyenthenguyen/docx v0.0.0-20230621112118-9c8e795a11db
Existing: github.com/ledongthuc/pdf v0.0.0-20220302134840-0c2507a12d80
2. Database Model
Created: 
document.go

New 
Document
 model with fields:

DocumentID - Primary key
UserID - Owner of the document
FileName - Original filename
FilePath - Storage path on disk
FileType - pdf, doc, or docx
FileSize - Size in bytes
ExtractedText - Full text content extracted from document
Summary
 - Optional AI-generated summary
CreatedAt, UpdatedAt - Timestamps
3. Database Migration
Created: 
create_documents_table.sql

SQL migration with:

Table structure for PostgreSQL (SERIAL primary key)
Indexes on userid and createdat for performance
Comments for documentation
IMPORTANT

Migration Required: Run the SQL migration to create the documents table before using the document component.

4. Document Processing Component
Created: 
comp_document.go

Comprehensive component with 4 actions:

Action: upload_and_extract
Accepts file upload via multipart form
Validates file type (PDF, DOC, DOCX only)
Enforces file size limit (10MB default, configurable)
Extracts text content using appropriate library
Stores metadata and extracted text in database
Returns document ID and extracted text
Injects data into ctx.Extras for downstream components
Action: list
Retrieves all documents for current user
Returns metadata without full extracted text (for performance)
Ordered by creation date (newest first)
Action: 
get
Retrieves specific document by ID
Security: Only returns documents owned by current user
Returns full metadata including extracted text
Injects document data into ctx.Extras
Action: delete
Deletes document from database and disk
Security: Only allows deletion of own documents
Graceful handling if file already deleted from disk
Text Extraction Functions:

extractTextFromPDF()
 - Extracts text from all pages of PDF files
extractTextFromDOCX()
 - Extracts text from DOCX files
Error handling for corrupted or unreadable files
5. AI Component Enhancement
Modified: 
comp_ai.go

Added document context support:

New Parameters:

document_ids - Comma-separated list of document IDs to include as context
use_all_documents - Boolean flag to include all user's documents
Implementation:

Document retrieval in 
delegateToLLM()
 function
Security: Only loads documents owned by current user
Text truncation: Limits each document to 5000 characters to avoid token overflow
Context injection: Adds document content to LLM system prompt
Format: Clearly labeled with document ID, filename, type, and content
Example Context Injection:

Available Documents:
[Document ID: 1 - training_manual.pdf]
Type: pdf | Size: 45231 bytes
Content:
This is the content of the training manual...
---
Build Verification
✅ Build Status: Successful

No compilation errors
All dependencies resolved
Component auto-registered via 
init()
 function
Manual Testing Guide
Prerequisites
Run the SQL migration to create the documents table
Start the backend server: go run main.go
Ensure you have a valid user session (JWT token)
Test 1: Upload and Extract PDF
Create a test PDF with sample content like:

Employee Training Manual
Section 1: Company Policies
All employees must follow the code of conduct...
Workflow Configuration:

Create a new workflow in the frontend
Add comp_document node:
action: upload_and_extract
file_field: document
Add comp_table node (optional) to save document ID
Execute:

Upload the PDF file
Verify response contains document_id and extracted_text
Check database: SELECT * FROM documents ORDER BY createdat DESC LIMIT 1;
Verify extractedtext field contains readable content
Test 2: AI with Document Context
Workflow Configuration:

Create a new workflow
Add comp_ai node:
command: "What does the training manual say about company policies?"
document_ids: 1 (or the ID from Test 1)
provider: gemini (or your configured LLM)
token: Your API key
Execute:

Run the workflow
Verify AI response references content from the uploaded document
Check console logs for: [CompAI] Loading documents with IDs: [1]
Verify: [CompAI] Loaded 1 documents into context
Test 3: Upload DOCX File
Create a test DOCX with content like:

Product Specifications
Model: XYZ-2024
Features: Advanced AI integration...
Execute:

Use same workflow as Test 1
Upload the DOCX file
Verify text extraction works correctly
Compare with PDF extraction quality
Test 4: List Documents
Workflow Configuration:

Add comp_document node:
action: list
Execute:

Verify response contains array of documents
Check that text_length is present but not full extracted_text
Verify documents are ordered by creation date
Test 5: Use All Documents
Workflow Configuration:

Upload 2-3 different documents (mix of PDF and DOCX)
Create AI workflow:
command: "Summarize all my documents"
use_all_documents: true
Execute:

Verify AI has access to all documents
Check console: [CompAI] Loading all documents for user X
Verify AI response references multiple documents
Test 6: Error Handling
Test unsupported file type:

Try uploading a .txt or .jpg file
Verify error: "UNSUPPORTED_FILE_TYPE"
Test file size limit:

Try uploading a file > 10MB
Verify error: "FILE_TOO_LARGE"
Test unauthorized access:

Try to get/delete another user's document
Verify error: "DOCUMENT_NOT_FOUND" or "access denied"
Key Features
✅ Security

User isolation: Users can only access their own documents
File type validation: Only PDF and DOCX allowed
File size limits: Prevents abuse (10MB default)
✅ Performance

Text truncation: Limits document content to 5000 chars per doc
List endpoint: Excludes full text for faster responses
Database indexes: Optimized queries on userid and createdat
✅ Robustness

Error handling: Graceful failures with cleanup
File cleanup: Removes files if database insert fails
Logging: Comprehensive debug output
✅ Integration

Context injection: Document data available to downstream components
Workflow engine: Results properly stored in wfEngine
AI enhancement: Seamless integration with existing AI component
Next Steps
Run the migration to create the documents table
Test the workflows using the manual testing guide above
Frontend integration: Add UI components for document upload
Optional enhancements:
Add document summary generation using AI
Implement vector search for semantic document retrieval
Add support for more file types (e.g., TXT, RTF)
Implement document versioning
Files Created
models/document.go
migrations/create_documents_table.sql
core/generator/comp_document.go
Files Modified
core/generator/comp_ai.go
 - Added document context support
go.mod
 - Added docx library dependency