-- Migration: Create documents table for storing uploaded documents and extracted text
-- This table supports the comp_document component for document processing

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

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_documents_userid ON documents(userid);
CREATE INDEX IF NOT EXISTS idx_documents_createdat ON documents(createdat);

-- Comments for documentation
COMMENT ON TABLE documents IS 'Stores uploaded documents with extracted text content for AI processing';
COMMENT ON COLUMN documents.extractedtext IS 'Full text content extracted from PDF/DOC/DOCX files';
COMMENT ON COLUMN documents.summary IS 'Optional AI-generated summary of document content';
