package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"erp6-be-golang/core/helpers"

	"github.com/gofiber/fiber/v2"
)

func init() {
	RegisterComponent("upload", func(ctx *WorkflowContext) error {
		return handleUpload(ctx)
	})
}

func handleUpload(ctx *WorkflowContext) error {
	c := ctx.FiberCtx
	params := ctx.Params

	var (
		filefield   string
		destination = "public/uploads"

		rename      = true
	)

	// Parse parameters
	for _, p := range params {
		switch p.InputName {
		case "filefield":
			filefield = resolveValue(ctx, strings.TrimSpace(p.CompValue))
		case "destination":
			val := resolveValue(ctx, strings.TrimSpace(p.CompValue))
			if val != "" {
				destination = val
			}
		case "rename":
			if p.CompValue == "false" {
				rename = false
			}
		case "enable":
			if p.CompValue == "false" {
				// If disabled, just return success immediately
				helpers.SuccessResponse(c, "UPLOAD SKIPPED", nil)
				return nil
			}
		}
	}

	if filefield == "" {
		return fmt.Errorf("filefield parameter is required for upload component")
	}

	// Security check for destination
	// Ensure destination is within allowed paths or just basic sanity check
	// For now, let's just clean the path to prevent directory traversal
	destination = filepath.Clean(destination)
	if strings.Contains(destination, "..") {
		helpers.FailResponse(c, fiber.StatusBadRequest, "INVALID PATH", "Destination path cannot contain '..'")
		return fmt.Errorf("invalid destination path: %s", destination)
	}

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(destination, 0755); err != nil {
		helpers.FailResponse(c, fiber.StatusInternalServerError, "UPLOAD ERROR", "Failed to create directory")
		return err
	}

	// Get file from request
	file, err := c.FormFile(filefield)
	if err != nil {
		fmt.Printf("[CompUpload] Error getting file from field '%s': %v\n", filefield, err)
		helpers.FailResponse(c, fiber.StatusBadRequest, "UPLOAD ERROR", "File not found in field: "+filefield)
		return err
	}
	fmt.Printf("[CompUpload] Received file: %s (Size: %d)\n", file.Filename, file.Size)

	// Generate filename
	filename := file.Filename
	if rename {
		ext := filepath.Ext(filename)
		name := strings.TrimSuffix(filename, ext)
		timestamp := time.Now().UnixNano()
		filename = fmt.Sprintf("%s_%d%s", name, timestamp, ext)
	}

	savePath := filepath.Join(destination, filename)
	fmt.Printf("[CompUpload] Saving to: %s\n", savePath)

	// Save file to disk
	if err := c.SaveFile(file, savePath); err != nil {
		helpers.FailResponse(c, fiber.StatusInternalServerError, "UPLOAD ERROR", "Failed to save file")
		return err
	}

	// Prepare result
	// Normalized path for web usage (forward slashes)
	webPath := strings.ReplaceAll(savePath, "\\", "/")
	
	resultData := map[string]interface{}{
		"filename":  filename,
		"path":      webPath,
		"size":      file.Size,
		"extension": filepath.Ext(filename),
	}

	// Store path in Extras so comp_table can use it
	// We store it with the key equal to the filefield name
	// e.g. if filefield is 'profile_pic', we store 'profile_pic' = 'public/uploads/filename.jpg'
	if ctx.Extras == nil {
		ctx.Extras = make(map[string]interface{})
	}
	ctx.Extras[filefield] = webPath

	// Store additional metadata
	ctx.Extras[filefield+"_path"] = webPath
	ctx.Extras[filefield+"_filename"] = filename
	ctx.Extras[filefield+"_original_filename"] = file.Filename
	ctx.Extras[filefield+"_size"] = file.Size
	ctx.Extras[filefield+"_mimetype"] = file.Header.Get("Content-Type")
	ctx.Extras[filefield+"_extension"] = filepath.Ext(filename)
	
	fmt.Printf("[CompUpload] Metadata stored in Extras for key '%s'\n", filefield)

	// Also Update wfEngine result
	if wfEngine, ok := c.Locals("wfEngine").([]WorkflowEngine); ok {
		wfEngine = append(wfEngine, WorkflowEngine{DataInputNode: "", ResultNode: resultData})
		c.Locals("wfEngine", wfEngine)
	}

	helpers.SuccessResponse(c, "FILE UPLOADED", resultData)

	return nil
}

func resolveValue(ctx *WorkflowContext, rawVal string) string {
	return os.Expand(rawVal, func(key string) string {
		var val string
		c := ctx.FiberCtx

		// Priority 1: Check Extras from previous workflow nodes or conversation state
		if ctx.Extras != nil {
			if extraVal, exists := ctx.Extras[key]; exists {
				val = fmt.Sprint(extraVal)
			}
		}

		// Priority 2: Check Locals (JWT token data)
		if val == "" {
			if key == "userid" {
				if userID, ok := c.Locals("userid").(int); ok && userID != 0 {
					val = fmt.Sprint(userID)
				}
			} else if key == "username" {
				if username, ok := c.Locals("username").(string); ok && username != "" {
					val = username
				}
			}
		}

		// Priority 3: Check Query and FormValue
		if val == "" {
			val = c.Query(key)
			if val == "" {
				val = c.FormValue(key)
			}
		}

		if val != "" {
			fmt.Printf("[CompUpload] Resolved variable '$%s' -> '%s'\n", key, val)
			return val
		}

		fmt.Printf("[CompUpload] Variable '$%s' resolved to empty\n", key)
		return ""
	})
}
