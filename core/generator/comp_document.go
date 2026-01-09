package generator

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"erp6-be-golang/core/helpers"
	"erp6-be-golang/models"

	"github.com/gofiber/fiber/v2"
	"github.com/ledongthuc/pdf"
	"github.com/nguyenthenguyen/docx"
	"gorm.io/gorm"
)

func init() {
	RegisterComponent("document", func(ctx *WorkflowContext) error {
		return handleDocument(ctx)
	})
}

func handleDocument(ctx *WorkflowContext) error {
	c := ctx.FiberCtx
	params := ctx.Params

	var (
		action       string
		fileField    string
		documentID   string
		maxFileSize  int64 = 10 * 1024 * 1024 // 10MB default limit
	)

	// Parse parameters
	for _, p := range params {
		switch p.InputName {
		case "action":
			action = strings.ToLower(strings.TrimSpace(p.CompValue))
		case "file_field", "filefield":
			fileField = strings.TrimSpace(p.CompValue)
		case "document_id", "documentid":
			documentID = ResolveParam(c, strings.TrimSpace(p.CompValue))
		case "max_file_size", "maxfilesize":
			// Optional: allow custom file size limit
			if p.CompValue != "" {
				fmt.Sscanf(p.CompValue, "%d", &maxFileSize)
			}
		}
	}

	// Get current user ID
	userID, ok := c.Locals("userid").(int)
	if !ok || userID == 0 {
		helpers.FailResponse(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "User ID not found")
		return fmt.Errorf("user ID not found in context")
	}

	switch action {
	case "upload_and_extract":
		return handleUploadAndExtract(ctx, fileField, documentID, userID, maxFileSize)
	case "list":
		return handleListDocuments(ctx, userID)
	case "get":
		return handleGetDocument(ctx, documentID, userID)
	case "delete":
		return handleDeleteDocument(ctx, documentID, userID)
	default:
		helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_ACTION", "Action must be: upload_and_extract, list, get, or delete")
		return fmt.Errorf("invalid action: %s", action)
	}
}

func handleUploadAndExtract(ctx *WorkflowContext, fileField string, documentID string, userID int, maxFileSize int64) error {
	c := ctx.FiberCtx
	db := ctx.DB

	if fileField == "" {
		helpers.FailResponse(c, fiber.StatusBadRequest, "MISSING_PARAMETER", "file_field is required for upload_and_extract")
		return fmt.Errorf("file_field parameter is required")
	}

	// Check if this is an UPDATE operation
	isUpdate := documentID != "" && documentID != "0"
	var existingDocument models.Document
	
	if isUpdate {
		// Verify document exists and belongs to user
		result := db.Where("documentid = ? AND userid = ?", documentID, userID).First(&existingDocument)
		if result.Error != nil {
			helpers.FailResponse(c, fiber.StatusNotFound, "DOCUMENT_NOT_FOUND", "Document not found or access denied")
			return result.Error
		}
		fmt.Printf("[CompDocument] UPDATE mode: Replacing document ID %s\n", documentID)
	} else {
		fmt.Printf("[CompDocument] INSERT mode: Creating new document\n")
	}

	var fileData []byte
	var fileName string
	var fileSize int64
	var contentType string

	// Try to get file from multipart form first
	file, err := c.FormFile(fileField)
	if err == nil {
		// Multipart file upload
		fmt.Printf("[CompDocument] Processing multipart file upload\n")
		
		// Open the file
		fileHandle, err := file.Open()
		if err != nil {
			helpers.FailResponse(c, fiber.StatusInternalServerError, "FILE_READ_ERROR", "Failed to read uploaded file")
			return err
		}
		defer fileHandle.Close()

		// Read file data
		fileData, err = io.ReadAll(fileHandle)
		if err != nil {
			helpers.FailResponse(c, fiber.StatusInternalServerError, "FILE_READ_ERROR", "Failed to read file content")
			return err
		}

		fileName = file.Filename
		fileSize = file.Size
		contentType = file.Header.Get("Content-Type")
	} else {
		// Try to get base64 encoded file from form field
		fmt.Printf("[CompDocument] Multipart file not found, checking for base64 data\n")
		
		base64Data := c.FormValue(fileField)
		if base64Data == "" {
			helpers.FailResponse(c, fiber.StatusBadRequest, "FILE_NOT_FOUND", "No file uploaded in field: "+fileField)
			return fmt.Errorf("no file data found in field: %s", fileField)
		}

		// Parse base64 data (format: data:mime/type;base64,<data>)
		var base64Content string
		if strings.HasPrefix(base64Data, "data:") {
			// Extract MIME type and base64 content
			parts := strings.SplitN(base64Data, ",", 2)
			if len(parts) != 2 {
				helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_BASE64", "Invalid base64 data format")
				return fmt.Errorf("invalid base64 format")
			}

			// Extract MIME type from data:mime/type;base64
			mimeTypePart := parts[0]
			if strings.Contains(mimeTypePart, ":") && strings.Contains(mimeTypePart, ";") {
				contentType = strings.TrimPrefix(strings.Split(mimeTypePart, ";")[0], "data:")
			}

			base64Content = parts[1]
		} else {
			// Plain base64 without data URI prefix
			base64Content = base64Data
		}

		// Decode base64
		var decodeErr error
		fileData, decodeErr = base64.StdEncoding.DecodeString(base64Content)
		if decodeErr != nil {
			helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID_BASE64", "Failed to decode base64 data")
			return fmt.Errorf("base64 decode error: %w", decodeErr)
		}

		fileSize = int64(len(fileData))

		// Get filename from separate field (file_filename)
		fileName = c.FormValue(fileField + "_filename")
		if fileName == "" {
			// Generate default filename based on content type
			ext := ".bin"
			if contentType != "" {
				switch contentType {
				case "application/pdf":
					ext = ".pdf"
				case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
					ext = ".docx"
				case "application/msword":
					ext = ".doc"
				}
			}
			fileName = fmt.Sprintf("document_%d%s", time.Now().Unix(), ext)
		}

		fmt.Printf("[CompDocument] Decoded base64 file: %s (Size: %d bytes)\n", fileName, fileSize)
	}

	fmt.Printf("[CompDocument] Validating file: %s (Size: %d bytes)\n", fileName, fileSize)

	// Validation 1: Check if file is empty
	if fileSize == 0 {
		helpers.FailResponse(c, fiber.StatusBadRequest, "EMPTY_FILE", "File is empty (0 bytes). Please upload a valid document.")
		return fmt.Errorf("empty file uploaded")
	}

	// Validation 2: Check minimum file size (100 bytes)
	minFileSize := int64(100)
	if fileSize < minFileSize {
		helpers.FailResponse(c, fiber.StatusBadRequest, "FILE_TOO_SMALL", 
			fmt.Sprintf("File is too small (%s). Minimum size is %s.", 
				formatFileSize(fileSize), formatFileSize(minFileSize)))
		return fmt.Errorf("file too small: %d bytes", fileSize)
	}

	// Validation 3: Check maximum file size
	if fileSize > maxFileSize {
		helpers.FailResponse(c, fiber.StatusBadRequest, "FILE_TOO_LARGE", 
			fmt.Sprintf("File size (%s) exceeds maximum limit of %s. Please upload a smaller file.", 
				formatFileSize(fileSize), formatFileSize(maxFileSize)))
		return fmt.Errorf("file too large: %d bytes (max: %d)", fileSize, maxFileSize)
	}

	// Validation 4: Check file extension
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext == "" {
		helpers.FailResponse(c, fiber.StatusBadRequest, "NO_FILE_EXTENSION", 
			"File has no extension. Please upload a file with .pdf or .docx extension.")
		return fmt.Errorf("file has no extension: %s", fileName)
	}

	// Validation 5: Validate supported file types
	allowedExtensions := map[string]string{
		".pdf":  "pdf",
		".docx": "docx",
		".doc":  "docx", // Legacy DOC files (will attempt DOCX extraction)
	}

	fileType, isSupported := allowedExtensions[ext]
	if !isSupported {
		helpers.FailResponse(c, fiber.StatusBadRequest, "UNSUPPORTED_FILE_TYPE", 
			fmt.Sprintf("File type '%s' is not supported. Allowed types: PDF (.pdf), Word (.docx, .doc)", ext))
		return fmt.Errorf("unsupported file type: %s", ext)
	}

	// Validation 6: Validate MIME type (additional security check)
	allowedMimeTypes := map[string]bool{
		"application/pdf":                                                                                  true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document":                         true, // .docx
		"application/msword":                                                                               true, // .doc
		"application/octet-stream":                                                                         true, // Generic binary (some browsers use this)
	}

	if contentType != "" && !allowedMimeTypes[contentType] {
		fmt.Printf("[CompDocument] Warning: Unexpected MIME type '%s' for file '%s'\n", contentType, fileName)
		// Don't fail here, just log warning - some browsers send incorrect MIME types
	}

	fmt.Printf("[CompDocument] Validation passed: Type=%s, Extension=%s, Size=%s\n", 
		fileType, ext, formatFileSize(fileSize))

	// Create upload directory
	uploadDir := "public/documents"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		helpers.FailResponse(c, fiber.StatusInternalServerError, "DIRECTORY_ERROR", "Failed to create upload directory")
		return err
	}

	// Generate unique filename
	timestamp := time.Now().UnixNano()
	uniqueFilename := fmt.Sprintf("%d_%s", timestamp, fileName)
	savePath := filepath.Join(uploadDir, uniqueFilename)

	// Save file to disk (works for both multipart and base64)
	if err := os.WriteFile(savePath, fileData, 0644); err != nil {
		helpers.FailResponse(c, fiber.StatusInternalServerError, "SAVE_ERROR", "Failed to save file")
		return err
	}

	fmt.Printf("[CompDocument] File saved to: %s\n", savePath)

	// Extract text based on file type
	var extractedText string
	if fileType == "pdf" {
		extractedText, err = extractTextFromPDF(savePath)
	} else {
		extractedText, err = extractTextFromDOCX(savePath)
	}

	if err != nil {
		// Clean up file if extraction fails
		os.Remove(savePath)
		helpers.FailResponse(c, fiber.StatusInternalServerError, "EXTRACTION_ERROR", "Failed to extract text from document")
		return err
	}

	fmt.Printf("[CompDocument] Extracted %d characters of text\n", len(extractedText))

	// Store or update in database
	webPath := strings.ReplaceAll(savePath, "\\", "/")
	
	var result *gorm.DB
	var document models.Document
	
	if isUpdate {
		// UPDATE existing document
		// Delete old file first
		if existingDocument.FilePath != "" {
			oldPath := existingDocument.FilePath
			if err := os.Remove(oldPath); err != nil {
				fmt.Printf("[CompDocument] Warning: Failed to delete old file %s: %v\n", oldPath, err)
			} else {
				fmt.Printf("[CompDocument] Deleted old file: %s\n", oldPath)
			}
		}
		
		// Update database record
		result = db.Model(&models.Document{}).
			Where("documentid = ? AND userid = ?", documentID, userID).
			Updates(map[string]interface{}{
				"filename":      fileName,
				"filepath":      webPath,
				"filetype":      fileType,
				"filesize":      fileSize,
				"extractedtext": extractedText,
				"updatedat":     time.Now(),
			})
		
		if result.Error != nil {
			// Clean up new file if update fails
			os.Remove(savePath)
			helpers.FailResponse(c, fiber.StatusInternalServerError, "DATABASE_ERROR", "Failed to update document")
			return result.Error
		}
		
		// Reload document to get updated data
		db.Where("documentid = ?", documentID).First(&document)
		fmt.Printf("[CompDocument] Document updated with ID: %s\n", documentID)
		
	} else {
		// INSERT new document
		document = models.Document{
			UserID:        userID,
			FileName:      fileName,
			FilePath:      webPath,
			FileType:      fileType,
			FileSize:      fileSize,
			ExtractedText: extractedText,
		}

		result = db.Create(&document)
		if result.Error != nil {
			// Clean up file if database insert fails
			os.Remove(savePath)
			helpers.FailResponse(c, fiber.StatusInternalServerError, "DATABASE_ERROR", "Failed to save document metadata")
			return result.Error
		}

		fmt.Printf("[CompDocument] Document created with ID: %d\n", document.DocumentID)
	}

	// Prepare response
	responseData := map[string]interface{}{
		"document_id":    document.DocumentID,
		"filename":       document.FileName,
		"file_type":      document.FileType,
		"file_size":      document.FileSize,
		"extracted_text": extractedText,
		"text_length":    len(extractedText),
		"is_update":      isUpdate,
	}

	// Inject into context for downstream components
	if ctx.Extras == nil {
		ctx.Extras = make(map[string]interface{})
	}
	ctx.Extras["document_id"] = document.DocumentID
	ctx.Extras["document_text"] = extractedText
	ctx.Extras["document_filename"] = document.FileName

	// Update workflow engine result
	wfEngine := c.Locals("wfEngine").([]WorkflowEngine)
	wfEngine = append(wfEngine, WorkflowEngine{DataInputNode: "", ResultNode: responseData})
	c.Locals("wfEngine", wfEngine)

	helpers.SuccessResponse(c, "DOCUMENT_UPLOADED", responseData)
	return nil
}

func handleListDocuments(ctx *WorkflowContext, userID int) error {
	c := ctx.FiberCtx
	db := ctx.DB

	var documents []models.Document
	result := db.Where("userid = ?", userID).Order("createdat DESC").Find(&documents)
	if result.Error != nil {
		helpers.FailResponse(c, fiber.StatusInternalServerError, "DATABASE_ERROR", "Failed to retrieve documents")
		return result.Error
	}

	// Prepare response (exclude full extracted text for list view)
	documentList := make([]map[string]interface{}, len(documents))
	for i, doc := range documents {
		documentList[i] = map[string]interface{}{
			"document_id": doc.DocumentID,
			"filename":    doc.FileName,
			"file_type":   doc.FileType,
			"file_size":   doc.FileSize,
			"text_length": len(doc.ExtractedText),
			"created_at":  doc.CreatedAt,
			"updated_at":  doc.UpdatedAt,
		}
	}

	responseData := map[string]interface{}{
		"documents": documentList,
		"count":     len(documentList),
	}

	// Update workflow engine result
	wfEngine := c.Locals("wfEngine").([]WorkflowEngine)
	wfEngine = append(wfEngine, WorkflowEngine{DataInputNode: "", ResultNode: responseData})
	c.Locals("wfEngine", wfEngine)

	helpers.SuccessResponse(c, "DOCUMENTS_RETRIEVED", responseData)
	return nil
}

func handleGetDocument(ctx *WorkflowContext, documentID string, userID int) error {
	c := ctx.FiberCtx
	db := ctx.DB

	if documentID == "" {
		helpers.FailResponse(c, fiber.StatusBadRequest, "MISSING_PARAMETER", "document_id is required")
		return fmt.Errorf("document_id parameter is required")
	}

	var document models.Document
	result := db.Where("documentid = ? AND userid = ?", documentID, userID).First(&document)
	if result.Error != nil {
		helpers.FailResponse(c, fiber.StatusNotFound, "DOCUMENT_NOT_FOUND", "Document not found or access denied")
		return result.Error
	}

	responseData := map[string]interface{}{
		"document_id":    document.DocumentID,
		"filename":       document.FileName,
		"file_path":      document.FilePath,
		"file_type":      document.FileType,
		"file_size":      document.FileSize,
		"extracted_text": document.ExtractedText,
		"summary":        document.Summary,
		"created_at":     document.CreatedAt,
		"updated_at":     document.UpdatedAt,
	}

	// Inject into context
	if ctx.Extras == nil {
		ctx.Extras = make(map[string]interface{})
	}
	ctx.Extras["document_id"] = document.DocumentID
	ctx.Extras["document_text"] = document.ExtractedText
	ctx.Extras["document_filename"] = document.FileName

	// Update workflow engine result
	wfEngine := c.Locals("wfEngine").([]WorkflowEngine)
	wfEngine = append(wfEngine, WorkflowEngine{DataInputNode: "", ResultNode: responseData})
	c.Locals("wfEngine", wfEngine)

	helpers.SuccessResponse(c, "DOCUMENT_RETRIEVED", responseData)
	return nil
}

func handleDeleteDocument(ctx *WorkflowContext, documentID string, userID int) error {
	c := ctx.FiberCtx
	db := ctx.DB

	if documentID == "" {
		helpers.FailResponse(c, fiber.StatusBadRequest, "MISSING_PARAMETER", "document_id is required")
		return fmt.Errorf("document_id parameter is required")
	}

	// First, get the document to retrieve file path
	var document models.Document
	result := db.Where("documentid = ? AND userid = ?", documentID, userID).First(&document)
	if result.Error != nil {
		helpers.FailResponse(c, fiber.StatusNotFound, "DOCUMENT_NOT_FOUND", "Document not found or access denied")
		return result.Error
	}

	// Delete file from disk
	if err := os.Remove(document.FilePath); err != nil {
		fmt.Printf("[CompDocument] Warning: Failed to delete file %s: %v\n", document.FilePath, err)
		// Continue with database deletion even if file deletion fails
	}

	// Delete from database
	result = db.Delete(&document)
	if result.Error != nil {
		helpers.FailResponse(c, fiber.StatusInternalServerError, "DATABASE_ERROR", "Failed to delete document")
		return result.Error
	}

	responseData := map[string]interface{}{
		"document_id": document.DocumentID,
		"filename":    document.FileName,
		"deleted":     true,
	}

	// Update workflow engine result
	wfEngine := c.Locals("wfEngine").([]WorkflowEngine)
	wfEngine = append(wfEngine, WorkflowEngine{DataInputNode: "", ResultNode: responseData})
	c.Locals("wfEngine", wfEngine)

	helpers.SuccessResponse(c, "DOCUMENT_DELETED", responseData)
	return nil
}

// extractTextFromPDF extracts text content from a PDF file
func extractTextFromPDF(filePath string) (string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	var textBuilder strings.Builder
	totalPages := r.NumPage()

	for pageNum := 1; pageNum <= totalPages; pageNum++ {
		page := r.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			fmt.Printf("[CompDocument] Warning: Failed to extract text from page %d: %v\n", pageNum, err)
			continue
		}

		textBuilder.WriteString(text)
		textBuilder.WriteString("\n\n")
	}

	return strings.TrimSpace(textBuilder.String()), nil
}

// extractTextFromDOCX extracts text content from a DOCX file
func extractTextFromDOCX(filePath string) (string, error) {
	r, err := docx.ReadDocxFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open DOCX: %w", err)
	}
	defer r.Close()

	doc := r.Editable()
	rawContent := doc.GetContent()

	// The GetContent() returns XML, we need to extract text from <w:t> tags
	var textBuilder strings.Builder
	
	// Simple XML parsing to extract text from <w:t> tags
	lines := strings.Split(rawContent, "<w:t")
	for i, line := range lines {
		if i == 0 {
			continue // Skip first part before any <w:t> tag
		}
		
		// Find the closing tag
		endTag := strings.Index(line, "</w:t>")
		if endTag == -1 {
			continue
		}
		
		// Find the end of opening tag
		startContent := strings.Index(line, ">")
		if startContent == -1 || startContent >= endTag {
			continue
		}
		
		// Extract text between tags
		text := line[startContent+1 : endTag]
		
		// Decode XML entities
		text = strings.ReplaceAll(text, "&lt;", "<")
		text = strings.ReplaceAll(text, "&gt;", ">")
		text = strings.ReplaceAll(text, "&amp;", "&")
		text = strings.ReplaceAll(text, "&quot;", "\"")
		text = strings.ReplaceAll(text, "&apos;", "'")
		
		textBuilder.WriteString(text)
	}
	
	// Add paragraph breaks by detecting <w:p> tags
	result := textBuilder.String()
	
	// Clean up: replace multiple spaces with single space
	result = strings.Join(strings.Fields(result), " ")
	
	// Add some basic paragraph structure
	paragraphs := strings.Split(rawContent, "<w:p ")
	var finalText strings.Builder
	
	for i, para := range paragraphs {
		if i == 0 {
			continue
		}
		
		// Extract text from this paragraph
		var paraText strings.Builder
		textParts := strings.Split(para, "<w:t")
		
		for j, part := range textParts {
			if j == 0 {
				continue
			}
			
			endTag := strings.Index(part, "</w:t>")
			if endTag == -1 {
				continue
			}
			
			startContent := strings.Index(part, ">")
			if startContent == -1 || startContent >= endTag {
				continue
			}
			
			text := part[startContent+1 : endTag]
			
			// Decode XML entities
			text = strings.ReplaceAll(text, "&lt;", "<")
			text = strings.ReplaceAll(text, "&gt;", ">")
			text = strings.ReplaceAll(text, "&amp;", "&")
			text = strings.ReplaceAll(text, "&quot;", "\"")
			text = strings.ReplaceAll(text, "&apos;", "'")
			
			paraText.WriteString(text)
		}
		
		if paraText.Len() > 0 {
			finalText.WriteString(paraText.String())
			finalText.WriteString("\n")
		}
	}

	return strings.TrimSpace(finalText.String()), nil
}

// extractTextFromDOC extracts text from legacy DOC files
// Note: This is a placeholder - proper DOC support requires additional libraries
func extractTextFromDOC(filePath string) (string, error) {
	// For now, return an error suggesting conversion
	return "", fmt.Errorf("legacy .doc format not supported - please convert to .docx or .pdf")
}

// formatFileSize converts bytes to human-readable format
func formatFileSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}
